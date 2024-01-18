package base

import "testing"

func TestMetrics(t *testing.T) {
	puzzle := NewEmptySudoku()

	if want, got := puzzle.MaxValueCount(), puzzle.ValueCount(); got != want {
		t.Errorf("ValueCount for an empty puzzle %d should equal MaxValueCount  %d", got, want)
	}

	given := func(x int, y int, value int) {
		puzzle.Cell(x, y).MustBe(value, Given, nil)
	}

	given(1, 1, 1)
	given(2, 2, 2)
	given(3, 3, 3)
	given(4, 4, 4)
	given(5, 5, 5)
	given(6, 6, 6)
	given(7, 7, 7)
	given(8, 8, 8)
	given(9, 9, 9)

	if want, got := puzzle.MaxValueCount() - 9 * (puzzle.Universe.Len() - 1), puzzle.ValueCount(); got != want {
		t.Errorf("After 9 givens, value count should be %d, but is %d", want, got)
	}

	// Force solve the puzzle
	for {
		if err := puzzle.DoConstraints(); err != nil {
			t.Fatalf("Error during DoConstraints: %s", err.Error())
		}
		if puzzle.IsSolved() {
			break
		}
		// Find a cell that isn't solved and pick a possible
		// value for it:
		for _, cell := range puzzle.Grid {
			if s, _ := cell.IsSolved(); !s {
				cell.Possibilities.DoValues(
					func(value int) bool {
						given(cell.X, cell.Y, value)
						return false})
				break
			}
		}
	}
	
	if want, got := puzzle.SolvedValueCount(), puzzle.ValueCount(); got != want {
		t.Errorf("ValueCount %d doesn't equal SolvedValueCount %d", got, want)
	}
}
