package html

import "testing"
import "sudoku/text"


var puzzle1 = `
123------
---456---
------789
987------
---321---
------654
4--------
-5-------
--6------
`  // end puzzle1

func TestHTML1(t *testing.T) {
	puzzle, err := text.TextToSudoku(puzzle1)
	if err != nil {
		t.Fatalf("Error parsing puzzle1: %s", err)
	}
	if err = puzzle.DoConstraints(); err != nil {
		t.Fatalf("Error during DoConstraints: %s", err.Error())
	}
	glyphs := make(map[int]rune)
	glyphs[5] = '你'
	table := ToTable(puzzle, glyphs)
	t.Logf("%s", table)
	t.Errorf("Look at HTML in log")
}


