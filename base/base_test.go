package base

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

	p.Cell(1, 1).MustBe(1, &Given{})
	p.Cell(1, 2).MustBe(2, &Given{})

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
	if _, err := p.Cell(1, 1).MustBe(2, &Given{}); err == nil {
		t.Errorf("Expected a contradiction")
	}
}
