package base

import "bytes"
import "testing"

func TestGridSetup(t *testing.T) {
	p := &Puzzle{}
	p.MakeCells(9)

	if p.Universe != 0x1FF {
		t.Errorf("Universe is %v", p.Universe)
	}

	c := p.Cell(3, 4)
	if c.X != 3 {
		t.Errorf("Cell at [3, 4] has X of %d", c.X)
	}
	if c.Y != 4 {
		t.Errorf("Cell at [3, 4] has Y of %d", c.Y)
	}
}

func TestRowColumn(t *testing.T) {
	p := &Puzzle{}
	p.MakeCells(9)
	p.AddLineGroups()
	p.Add3x3Groups()

	p.Cell(1, 1).MustBe(1, Given)
	p.Cell(1, 2).MustBe(2, Given)

/*
	if b, v := p.Cell(1, 1).IsSolved(); b {
		if v != 1 {
			t.Errorf("IsSolved returned wring value")
		}
	} else {
		t.Errorf("IsSolved should have returned true: %d", p.Cell(1, 1).Possibilities)
	}
*/
	if err := p.DoConstraints(); err != nil {
		t.Errorf("Error during DoConstraints: %s", err.Error())
	}

	checkCell := func(x int, y int, value int, expect bool) {
		c := p.Cell(x, y)
		if got := c.HasPossibleValue(value); got != expect {
			t.Errorf("Cell(%d, %d)HasValue(%d) possibilities: %x: got %v, want %v", x, y, value, c.Possibilities, got, expect)
		}
	}
	checkCell(1, 1, 1, true)
	checkCell(1, 1, 9, false)
	checkCell(1, 2, 1, false)
	checkCell(1, 3, 1, false)
	checkCell(1, 3, 9, true)
	checkCell(3, 3, 1, false)
	checkCell(3, 3, 2, false)
	checkCell(3, 3, 4, true)
	if _, err := p.Cell(1, 1).MustBe(2, Given); err == nil {
		t.Errorf("Expected a contradiction")
	}
}

func TestSudoku1(t *testing.T) {
	// This test is from a sudoku published on page 10 of MIT's
	// student newspaper The Tech on 2017-09-21.
	p := &Puzzle{}
	p.MakeCells(9)
	p.AddLineGroups()
	p.Add3x3Groups()
	given := func(x int, y int, value int) {
		p.Cell(x, y).MustBe(value, Given)	
	}
	given(3, 1, 2)
	given(4, 1, 8)
	given(1, 2, 7)
	given(3, 2, 5)
	given(5, 2, 1)
	given(7, 2, 8)
	given(8, 2, 2)
	given(2, 3, 8)
	given(4, 3, 7)
	given(8, 3, 1)
	given(4, 4, 4)
	given(7, 4, 5)
	given(8, 4, 7)
	given(9, 4, 3)
	given(4, 5, 1)
	given(5, 5, 6)
	given(6, 5, 7)
	given(1, 6, 2)
	given(2, 6, 7)
	given(3, 6, 4)
	given(6, 6, 5)
	given(2, 7, 6)
	given(6, 7, 8)
	given(8, 7, 5)
	given(2, 8, 2)
	given(3, 8, 8)
	given(5, 8, 7)
	given(7, 8, 9)
	given(9, 8, 6)
	given(6, 9, 3)
	given(7, 9, 4)

	// Must use go test -v to see the log output of a successful test.
	show := func() {
		var b bytes.Buffer
		p.Show(&b)
		t.Log(b.String())
	}

	show()

	if err := p.DoConstraints(); err != nil {
		t.Errorf("Error during DoConstraints: %s", err.Error())
	}
	show()
	for _, j := range p.Justifications {
		t.Log(j.Pretty())
	}
}

func TestSudoku2(t *testing.T) {
	// This test is from the "Pset Nights" sudoku published on page 9 of MIT's
	// student newspaper The Tech on 2017-09-14.
	p := &Puzzle{}
	p.MakeCells(9)
	p.AddLineGroups()
	p.Add3x3Groups()
	given := func(x int, y int, value int) {
		p.Cell(x, y).MustBe(value, Given)	
	}
	given(2, 1, 3)
	given(5, 1, 2)
	given(7, 1, 6)
	given(8, 1, 7)
	given(9, 1, 4)
	given(1, 2, 2)
	given(1, 3, 7)
	given(2, 3, 6)
	given(3, 3, 8)
	given(6, 3, 5)
	given(8, 3, 2)
	given(4, 4, 2)
	given(9, 4, 5)
	given(2, 5, 2)
	given(4, 5, 9)
	given(6, 5, 4)
	given(8, 5, 3)
	given(1, 6, 4)
	given(6, 6, 1)
	given(2, 7, 4)
	given(4, 7, 6)
	given(7, 7, 2)
	given(8, 7, 9)
	given(9, 7, 3)
	given(9, 8, 6)
	given(1, 9, 9)
	given(2, 9, 9)
	given(3, 9, 6)
	given(5, 9, 3)
	given(8, 9, 4)

	show := func() {
		var b bytes.Buffer
		p.Show(&b)
		t.Log(b.String())
	}

	show()

	if err := p.DoConstraints(); err != nil {
		t.Errorf("Error during DoConstraints: %s", err.Error())
	}
	show()
	for _, j := range p.Justifications {
		t.Log(j.Pretty())
	}
}
