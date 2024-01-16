// Package text provides a way to set up Sudoku and KenKen puzzles from a
// textual representation.
package text

import "bufio"
import "fmt"
import "io"
import "regexp"
import "strconv"
import "strings"
import "unicode"
import "unicode/utf8"
import "sudoku/base"

// TextToSudoku returns an unsolved puzzle representing the specified sudoku.
// The string argument should be a string of digits representing the given
// values and dashes representing empty cells.  Spaces and tabs are ignored.
// Newlines represent breaks between rows.
func TextToSudoku(text string) (*base.Puzzle, error) {
	p := base.NewEmptySudoku()

	linenumber := 1
	linecharnumber := 0
	row := 1
	column := 1
	any := false
	comment := false

	new_line := func () {
		if any {
			row += 1
			column = 1
		}
		linenumber += 1
		linecharnumber = 0
	}

	for _, c := range text {
		linecharnumber += 1
		if comment {
			if c == '\n' {
				comment = false
				new_line()
			}
			continue
		}
		switch c {
		case '#':
			comment = true
			break
		case ' ', '\t':
			// Ignore
			break
		case '\n':
			new_line()
			break
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if row > 9 || column > 9 {
				return p, fmt.Errorf("row or column index overflow: row %d, column %d",
					row, column)
			}
			value := int(c - '0')
			p.Cell(column, row).MustBe(value, base.Given, nil)
			fallthrough
		case '-':
			any = true
			column += 1
			break
		default:
			return p, fmt.Errorf("invalid input character '%c': line %d, character %d",
				c, linenumber, linecharnumber)
		}
	}

	return p, nil
}

func max(ints ...int) int {
	m := ints[0]
	for _, i := range ints[1:] {
		if i > m {
			m = i
		}
	}
	return m
}

var CageConstraintRegexp = regexp.MustCompile(
	"[ \t]*(?P<group>[a-zA-Z])[ \t]*:[ \t]*(?P<value>[0-9]+)[ \t]*(?P<op>[-+*/])")

// TextToKenKen makes a ken-ken puzzle from a text specification.
// The specification starts with a grid of letters, digits and hyphens.
// Cells identified by the same letter are in the same ken-ken cage.
// Cells marked with a digit contain that fixed value.
// Cells marked with a hyphen are not in any cage.
// After the grid description are the rules for each cage identifying
// the operator and resulting value.
func TextToKenKen(text string) (*base.Puzzle, error) {
	p := &base.Puzzle{
		Grid: make(map[base.GridKey]*base.Cell),
		// Groups: make([]*Group, 0),
	}
	groups := make(map[rune]*base.Group)
	cell_values := make(map[*base.Cell]int)
	size := 0
	last_cell_row := 0

	cell := func(x, y int) *base.Cell {
		key := base.MakeGridKey(x, y)
		c := p.Grid[key]
		if c == nil {
			c = &base.Cell{
				X:      x,
				Y:      y,
				Puzzle: p,
			}
			p.Grid[key] = c
		}
		size = max(size, x, y)
		last_cell_row = y
		return c
	}

	group := func(c rune) *base.Group {
		g := groups[c]
		if g == nil {
			g = base.NewGroup(p)
			groups[c] = g
			p.Groups = append(p.Groups, g)
		}
		return g
	}

	// First read the grid.
	linenumber := 1
	linecharnumber := 0
	row := 1
	column := 1
	any := false
	comment := false

	new_line := func () {
		if any {
			row += 1
			column = 1
		}
		linenumber += 1
		linecharnumber = 0
	}

	next_cell := func() {
		any = true
		column += 1
	}

	reader := bufio.NewReader(strings.NewReader(text))

	for {
		c, _, err := reader.ReadRune()
		if err != nil {
			return p, err
		}
		if comment {
			if c == '\n' {
				comment = false
				new_line()
			}
			continue
		}
		switch c {
		case '#':
			comment = true
			break
		case ' ', '\t':
			// Ignore
			break
		case '\n':
			if any {
				row += 1
				column = 1
			}
			// Two empty lines means the grid is done.
			if row-last_cell_row >= 2 {
				goto grid_done
			}

			break
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			value := int(c - '0')
			cell_values[cell(column, row)] = value
			size = max(size, value)
			next_cell()
			break
		case '-':
			cell(column, row)
			next_cell()
			break
		default:
			if !unicode.IsLetter(c) {
				return p, fmt.Errorf("invalid input character '%c': line %d, character %d",
					c, linenumber, linecharnumber)
			}
			cl := cell(column, row)
			g := group(c)
			g.AddCell(cl)
			next_cell()
			break
		}
	}
grid_done:

	p.Size = size
	p.Universe = base.Universe(size)
	for _, c := range p.Grid {
		c.Possibilities = p.Universe
	}

	p.AddLineGroups()

	for c, v := range cell_values {
		c.MustBe(v, base.Given, nil)
	}

	// Now read the cage constraints.
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return p, nil
			}
			return p, err
		}
		m := CageConstraintRegexp.FindStringSubmatch(s)
		if m == nil {
			continue
		}
		group_identifier := m[1]
		value_str := m[2]
		value, err := strconv.Atoi(value_str)
		if err != nil {
			return p, err
		}
		op_str := m[3]
		op := base.KenKenOperatorSymbols[op_str]
		if op == nil {
			return p, fmt.Errorf("unsupported cage operator symbol %s", op_str)
		}

		gi, _ := utf8.DecodeRuneInString(group_identifier)
		group := groups[gi]
		if group == nil {
			return p, fmt.Errorf("There's no group for copnstraint %s", group_identifier)
		}
		group.AddConstraint(base.MakeKenKenConstraint([]*base.KenKenOperator{op}, value))
	}

	return p, nil
}
