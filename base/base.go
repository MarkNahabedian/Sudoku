// Package base implements the basic data model and functionality for
// a sudoku solver.
package base

import "fmt"
import "io"
import "os"

// Contradiction is the type of error that is returned if an operation
// results in a contradiction.
type Contradiction struct {
	// Cell is the Cell that's the subject of the contradiction.
	Cell *Cell
	// Issue is the error message describing the contradiction.
	Issue string
}

// Error implements the error interface.
func (c *Contradiction) Error() string {
	return fmt.Sprintf("Contradiction at [%d, %d]: %s", c.Cell.X, c.Cell.Y, c.Issue)
}

// GridKey represents a location in a Puzzle's cell grid.
type GridKey struct {
	X int
	Y int
}

// MakeGridKey returns a GridKey for the specified X and Y position.
func MakeGridKey(x int, y int) GridKey {
	return GridKey{X: x, Y: y}
}

// Cell represents a single cell of a Puzzle.  When the Puzzle is solved,
// each of its Cells will contain a single value.
type Cell struct {
	X             int
	Y             int
	Puzzle        *Puzzle
	Groups        []*Group
	Possibilities ValueSet
}

// Puzzle represents a single sudoku puzzle.
type Puzzle struct {
	// Size is the numbner of rows and columns in the Puzzle grid.
	Size int
	// Grid is the grid of cells.
	Grid map[GridKey]*Cell
	// Progress is a monotonically increasing integer.
	// It increases whenever progress is made.
	Progress uint
	// Groups is a slice conbtaining all of the Groups of this Puzzle.
	Groups []*Group
	// Universe is a ValueSet of all of the values that could appear in any Cell.
	Universe ValueSet
	// Justifications is a slice of all of the Justifications for what's
	// been asserted about this Puzzle.
	Justifications []*Justification
}

// Cell returns the Cell at the specified x and y position.
func (p *Puzzle) Cell(x, y int) *Cell {
	c := p.Grid[MakeGridKey(x, y)]
	if c == nil {
		panic(fmt.Sprintf("No cell at %d, %d", x, y))
	}
	return c
}

func (p *Puzzle) MakeCells(size int) *Puzzle {
	p.Size = size
	for i := 1; i <= size; i++ {
		p.Universe = p.Universe.SetHasValue(i, true)
	}
	p.Grid = make(map[GridKey]*Cell)
	for x := 1; x <= size; x++ {
		for y := 1; y <= size; y++ {
			p.Grid[MakeGridKey(x, y)] = &Cell{
				Puzzle:        p,
				X:             x,
				Y:             y,
				Groups:        make([]*Group, 0),
				Possibilities: p.Universe,
			}
		}
	}
	return p
}

func (p *Puzzle) Show(f io.Writer) {
	fmt.Fprintf(f, "\n")
	for y := 1; y <= p.Size; y++ {
		for x := 1; x <= p.Size; x++ {
			c := p.Cell(x, y)
			if b, v := c.IsSolved(); b {
				fmt.Fprintf(f, "  %d  ", v)
			} else {
				fmt.Fprintf(f, " %03x ", c.Possibilities)
			}
		}
		fmt.Fprintf(f, "\n")
	}
	fmt.Fprintf(f, "\n")
}

func (p *Puzzle) AddGroup(g *Group) *Puzzle {
	if g.Puzzle() != p {
		panic("Group.Puzzle() doesn't match Puzzle")
	}
	p.Groups = append(p.Groups, g)
	for _, c := range g.Cells() {
		c.Groups = append(c.Groups, g)
	}
	return p
}

func (c *Cell) IsSolved() (bool, int) {
	value := -1
	count := 0
	c.Possibilities.DoValues(func(v int) bool {
		count += 1
		value = v
		return true
	})
	if count == 1 {
		return true, value
	} else {
		return false, value
	}
}

func (c *Cell) HasPossibleValue(v int) bool {
	return c.Possibilities.HasValue(v)
}

type JustificationOp int

const (
	MUST_BE JustificationOp = iota
	CANT_BE
)

var JustificationOpStrings map[JustificationOp]string = map[JustificationOp]string{
	MUST_BE: "MUST_BE",
	CANT_BE: "CANT_BE",
}

type Justification struct {
	Tick       uint
	Cell       *Cell
	Constraint Constraint
	Operation  JustificationOp
	Value      int
}

func (j *Justification) Pretty() string {
	return fmt.Sprintf("%3d: Cell(%d, %d) %s %d %s",
		j.Tick, j.Cell.X, j.Cell.Y, JustificationOpStrings[j.Operation],
		j.Value, j.Constraint.Name())
}

func (p *Puzzle) Justify(c *Cell, op JustificationOp, value int, constraint Constraint) *Justification {
	j := &Justification{
		Tick:       p.Progress,
		Cell:       c,
		Constraint: constraint,
		Operation:  op,
		Value:      value,
	}
	p.Progress += 1
	p.Justifications = append(p.Justifications, j)
	return j
}

func (c *Cell) CantBe(v int, constraint Constraint) *Cell {
	old := c.Possibilities
	c.Possibilities = c.Possibilities.SetHasValue(v, false)
	if c.Possibilities != old {
		c.Puzzle.Justify(c, CANT_BE, v, constraint)
		if c.Possibilities.IsEmpty() {
			fmt.Fprintf(os.Stderr, "No remaining possible values for cell %d, %d\n", c.X, c.Y)
			for _, j := range c.Puzzle.Justifications {
				if j.Cell == c {
					fmt.Fprintf(os.Stderr, "%d: %s %d value: %d\n",
						j.Tick, j.Constraint, j.Operation, j.Value)
				}
			}
		}
	}
	return c
}

func (c *Cell) MustBe(v int, constraint Constraint) (*Cell, error) {
	old := c.Possibilities
	if !c.HasPossibleValue(v) {
		return c, &Contradiction{Cell: c, Issue: fmt.Sprintf("%d is not a possible Value for MustBe", v)}
	}
	c.Possibilities = NewValueSet([]int{v})
	if c.Possibilities.IsEmpty() {
		panic("MustBe produced Empty")
	}
	if c.Possibilities != old {
		c.Puzzle.Justify(c, MUST_BE, v, constraint)
	}
	return c, nil
}

func (p *Puzzle) DoConstraints() error {
	for {
		startTick := p.Progress
		for _, g := range p.Groups {
			if err := g.DoConstraints(); err != nil {
				return err
			}
		}
		if p.Progress == startTick {
			break
		}
	}
	return nil
}

// Constraint implements a Constraint that is to be enforced on
// a Group of Cells.
type Constraint interface {
	Name() string
	DoConstraint(*Group) error
}

type FunctionConstraint struct {
	name       string
	constraint func(*Group) error
}

func (constraint FunctionConstraint) Name() string {
	return constraint.name
}

func (constraint FunctionConstraint) DoConstraint(g *Group) error {
	return constraint.constraint(g)
}

var Given = FunctionConstraint{
	name:       "Given",
	constraint: (func(g *Group) error { return nil }),
}

// Group represents a set of cells to which some constraint collectively
// applies, for example a row in a sudoku whose cells can not contain
// matching values.
type Group struct {
	puzzle      *Puzzle
	cells       []*Cell
	constraints []Constraint
}

func (g *Group) Puzzle() *Puzzle { return g.puzzle }

func (g *Group) Cells() []*Cell { return g.cells }

// DoConstraints applies the Group's constraints.
func (g *Group) DoConstraints() error {
	for _, constraint := range g.constraints {
		if err := constraint.DoConstraint(g); err != nil {
			return err
		}
	}
	return nil
}

// HereThenNotElsewhereConstraint implements the constraint that says that
// the value of one Cell of a Group can not appear in any other Cell of
// that Group.
var HereThenNotElsewhereConstraint FunctionConstraint

func init() {
	HereThenNotElsewhereConstraint.name = "HereThenNotElsewhereConstraint"
	HereThenNotElsewhereConstraint.constraint = func(g *Group) error {
		// If there are n cells in the group with n possibilities and those
		// possibilities are the same then none of the other cells in the group
		// can have any of those values.
		for i, c1 := range g.Cells() {
			if c1.Possibilities.Len() > c1.Puzzle.Universe.Len()-i {
				continue
			}
			count := 0 // How many of the remaining cells have the same Possibilities
			for _, c2 := range g.Cells()[i:] {
				if c1.Possibilities == c2.Possibilities {
					count += 1
				}
			}
			if count == c1.Possibilities.Len() {
				for _, c3 := range g.Cells() {
					if c3.Possibilities != c1.Possibilities {
						c1.Possibilities.DoValues(func(val int) bool {
							c3.CantBe(val, HereThenNotElsewhereConstraint)
							return true
						})
					}
				}
			}
		}
		return nil
	}
}

// NotElsewhereThenHereConstraint implements the constraint that says
// that each value must appear in some Cell of the Group.
var NotElsewhereThenHereConstraint FunctionConstraint

func init() {
	NotElsewhereThenHereConstraint.name = "NotElsewhereThenHereConstraint"
	NotElsewhereThenHereConstraint.constraint = func(g *Group) error {
		// If a value can't be in all but one cell of a group then the one
		// cell that can have that value must have it.
		// Is there a generalization of this for several values?
		valueCells := make(map[int][]*Cell)
		for _, c := range g.Cells() {
			c.Possibilities.DoValues(func(v int) bool {
				valueCells[v] = append(valueCells[v], c)
				return true
			})
		}
		for v, cells := range valueCells {
			if len(cells) == 1 {
				if _, err := cells[0].MustBe(v, NotElsewhereThenHereConstraint); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// AddLinerGroups adds the conventional sudoku row and column constraints
// to the Puzzle.
func (p *Puzzle) AddLineGroups() *Puzzle {
	addGroup := func(cells []*Cell) {
		g := &Group{
			puzzle: p,
			cells:  cells,
			constraints: []Constraint{
				HereThenNotElsewhereConstraint,
				NotElsewhereThenHereConstraint,
			},
		}
		p.AddGroup(g)
	}
	for x := 1; x <= p.Size; x++ {
		column := []*Cell{}
		for y := 1; y <= p.Size; y++ {
			column = append(column, p.Cell(x, y))
		}
		addGroup(column)
	}
	for y := 1; y <= p.Size; y++ {
		row := []*Cell{}
		for x := 1; x <= p.Size; x++ {
			row = append(row, p.Cell(x, y))
		}
		addGroup(row)
	}
	return p
}

// Add3x3Groups implements the small 3x3 box constraints of a sudoku.
func (p *Puzzle) Add3x3Groups() *Puzzle {
	for sx := 1; sx < 9; sx += 3 {
		for sy := 1; sy < 9; sy += 3 {
			block := []*Cell{}
			for dx := 0; dx < 3; dx++ {
				for dy := 0; dy < 3; dy++ {
					x := sx + dx
					y := sy + dy
					block = append(block, p.Cell(x, y))
				}
			}
			g := &Group{
				puzzle: p,
				cells:  block,
				constraints: []Constraint{
					HereThenNotElsewhereConstraint,
					NotElsewhereThenHereConstraint,
				},
			}
			p.AddGroup(g)
		}
	}
	return p
}

type KenKenOperator struct {
	// Symbol is the name of the operator.
	Symbol string
	// Test returns true if the constraint is satisfied.
	Test func([]int, int) bool
}

var KenKenOperators []KenKenOperator = []KenKenOperator{
	{
		Symbol: "Addition",
		// OrderMatters: false,
		/*
			Function: func(values []int) int {
				sum := 0
				for _, v := range values {
					sum += v
				}
				return sum
			}
		*/
		Test: func(values []int, expect int) bool {
			sum := 0
			for _, v := range values {
				sum += v
			}
			return sum == expect
		},
	},
	{
		Symbol: "Multiplication",
		Test: func(values []int, expect int) bool {
			product := 1
			for _, v := range values {
				product *= v
			}
			return product == expect
		},
	},
	{
		Symbol: "Subtraction",
		Test: func(values []int, expect int) bool {
			for sense := 0; sense < 2^len(values); sense++ {
				accumulator := 0
				for index, value := range values {
					if sense&(1<<uint(index)) == 0 {
						accumulator += value
					} else {
						accumulator += -value
					}
				}
				if accumulator == expect {
					return true
				}
			}
			return false
		},
	},
}

// GetKenKenOperator returns the KenKenOperator with the specified name,
// or nil if none is found.
func GetKenKenOperator(name string) *KenKenOperator {
	for _, o := range KenKenOperators {
		if o.Symbol == name {
			return &o
		}
	}
	return nil
}

type KenKenCageConstraint struct {
	name      string
	operators []*KenKenOperator
	expect    int
}

func (c *KenKenCageConstraint) makeName() string {
	opString := ""
	for _, op := range c.operators {
		if opString != "" {
			opString += ", "
		}
		opString += op.Symbol
	}
	return opString
}

func (c *KenKenCageConstraint) Name() string {
	if c.name == "" {
		c.name = c.makeName()
	}
	return c.name
}

func (c *KenKenCageConstraint) findPossibleValueSets() {
	// possibilities := [][]int{}
	// Need group size
}

func (c *KenKenCageConstraint) DoConstraint(g *Group) error {
	/*
	   We might be able to do union and intersection of value sets when updating cells.
	*/

	cell_count := len(g.cells)
	cell_value_indices := make([]int, cell_count, cell_count)
	successful_possibilities := [][]int{}
	done := false
	for !done {
		// Test the current cell values
		values := make([]int, cell_count, cell_count)
		for i := 0; i < cell_count; i++ {
			values[i] = g.cells[i].Possibilities.Get(cell_value_indices[i])
		}
		for _, o := range c.operators {
			if o.Test(values, c.expect) {
				successful_possibilities = append(successful_possibilities, values)
			}
		}
		// Next cell value
		for cell_index := 0; true; cell_index++ {
			if cell_index >= cell_count {
				done = true
				break
			}
			cell_value_indices[cell_index] += 1
			if cell_value_indices[cell_index] < g.Cells()[cell_index].Possibilities.Len() {
				break
			}
			cell_value_indices[cell_index] = 0
		}
	}
	// Eliminate each cell value possibility that does not appear in
	// any of the acceptable combinations.
	for cell_index := 0; cell_index < cell_count; cell_index++ {
		cell := g.Cells()[cell_index]
		cell.Possibilities.DoValues(func(p int) bool {
			found := false
			for _, sp := range successful_possibilities {
				if sp[cell_index] == p {
					found = true
					break
				}
			}
			if !found {
				cell.CantBe(p, c)
			}
			return true
		})
	}
	return nil
}

func MakeKenKenConstraint(operators []*KenKenOperator, expect int) Constraint {
	return &KenKenCageConstraint{
		operators: operators,
		expect:    expect,
	}
}
