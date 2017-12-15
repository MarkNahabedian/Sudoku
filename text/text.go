// Package text provides a way to set up Sudoku and KenKen puzzles from a
// textual representation.
package text

import "fmt"
import "sudoku/base"

// TextToSudoku returns an unsolved puzzle representing the specified sudoku.
// The string argument should be a string of digits representing the given
// values and dashes representing empty cells.  Spaces and tabs are ignored.
// Newlines represent breaks between rows.
func TextToSudoku(text string) (*base.Puzzle, error) {
	p := &base.Puzzle{}
	p.MakeCells(9)
	p.AddLineGroups()
	p.Add3x3Groups()

	row := 1
	column := 1
    any := false
	for index, c := range text {
		switch c {
		case ' ','\t':
			// Ignore
            break
		case '\n':
        	if any {
				row += 1
				column = 1
            }
            break
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
        	if row > 9 || column > 9 {
				return p, fmt.Errorf("row or column index overflow at character %d", index)
            }
			value := int(c - '0')
			p.Cell(column, row).MustBe(value, base.Given)
            fallthrough
		case '-':
			any = true
			column += 1
            break
		}
	}

	return p, nil
}

