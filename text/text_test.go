package text

import "testing"

import "bytes"

func TestSudokuOnly17Given(t *testing.T) {
	p, err := TextToSudoku(`
		---7-----
		1--------
		---43-2--
		--------6
		---5-9---
		------418
		----81---
		--2----5-
		-4----3--
	`)

    if err != nil {
		t.Errorf("%s", err)
	}

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

