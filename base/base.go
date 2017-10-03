package base

import "fmt"
import "io"
import "os"
import "reflect"
import "runtime"

type Contradiction struct {
	Cell *Cell
	Issue string
}

func (c *Contradiction) Error() string {
	return fmt.Sprintf("Contradiction at [%d, %d]: %s", c.Cell.X, c.Cell.Y, c.Issue)
}

type GridKey struct {
	X int
	Y int
}

func MakeGridKey(x int, y int) GridKey {
	return GridKey{ X: x, Y: y }
}

type Cell struct {
	X int
	Y int
	Puzzle *Puzzle
	Groups []*Group
	Possibilities ValueSet
}

type Puzzle struct {
	Size int
	Grid map[GridKey]*Cell
	Progress uint  // should be increased whenever progrress is made
	Groups []*Group
	Universe ValueSet
	Justifications []*Justification
}

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
				Puzzle: p,
				X: x,
				Y: y,
				Groups: make([]*Group, 0),
				Possibilities: p.Universe,
			}
		}
	}
	return p
}

func (p *Puzzle) Show(f io.Writer) {
	fmt.Fprintf(f, "\n")
	for y:= 1; y <= p.Size; y++ {
		for x:= 1; x <= p.Size; x++ {
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
	c.Possibilities.DoValues(func (v int) {
		count += 1
		value = v
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

var JustificationOpStrings map[JustificationOp]string = map[JustificationOp]string {
	MUST_BE: "MUST_BE",
	CANT_BE: "CANT_BE",
}

type Justification struct {
	Tick uint
	Cell *Cell
	Constraint Constraint
	Operation JustificationOp
	Value int
}

func (j *Justification) Pretty() string {
	return fmt.Sprintf("%3d: Cell(%d, %d) %s %d %s",
		j.Tick, j.Cell.X, j.Cell.Y, JustificationOpStrings[j.Operation],
		j.Value, j.Constraint.Name())
}

func (p *Puzzle) Justify(c *Cell, op JustificationOp, value int, constraint Constraint) *Justification {
	j := &Justification {
		Tick: p.Progress,
		Cell: c,
		Constraint: constraint,
		Operation: op,
		Value: value,
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
		return c, &Contradiction{ Cell: c, Issue: fmt.Sprintf("%d is not a possible Value for MustBe", v) }
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

type Constraint func(*Group) error

func (constraint Constraint) Name() string {
	return runtime.FuncForPC(reflect.ValueOf(constraint).Pointer()).Name()
}

func Given(g *Group) error { return nil }

type Group struct {
	puzzle *Puzzle
	cells []*Cell
	constraints []Constraint
}

func (g *Group) Puzzle() *Puzzle { return g.puzzle }

func (g *Group) Cells() []*Cell { return g.cells }

func (g *Group) DoConstraints() error {
	for _, constraint := range g.constraints {
		if err := constraint(g); err != nil {
			return err
		}
	}
	return nil
}

func HereThenNotElsewhereConstraint(g *Group) error {
	// If there are n cells in the group with n possibilities and those
	// possibilities are the same then none of the other cells in the group
	// can have any of those values.
	for i, c1 := range g.Cells() {
		if c1.Possibilities.Len() > c1.Puzzle.Universe.Len() - i {
			continue
		}
		count := 0  // How many of the remaining cells have the same Possibilities
		for _, c2 := range g.Cells()[i:] {
			if c1.Possibilities == c2.Possibilities {
				count += 1
			}
		}
		if count == c1.Possibilities.Len() {
			for _, c3 := range g.Cells() {
				if c3.Possibilities != c1.Possibilities {
					c1.Possibilities.DoValues(func(val int) {
						c3.CantBe(val, HereThenNotElsewhereConstraint)
					})
				}
			}
		}
	}
	return nil
}

func NotElsewhereThenHereConstraint(g *Group) error {
	// If a value can't be in all but one cell of a group then the one
	// cell that can have that value must have it.
	// Is there a generalization of this for several values?
	valueCells := make(map[int][]*Cell)
	for _, c := range g.Cells() {
		c.Possibilities.DoValues(func (v int) {
			valueCells[v] = append(valueCells[v], c)
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

func (p *Puzzle) AddLineGroups() *Puzzle {
	addGroup := func(cells []*Cell) {
		g := &Group{
			puzzle: p,
			cells: cells,
			constraints: []Constraint {
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
				cells: block,
				constraints: []Constraint {
					HereThenNotElsewhereConstraint,
					NotElsewhereThenHereConstraint,
				},
			}
			p.AddGroup(g)
		}
	}
	return p
}
