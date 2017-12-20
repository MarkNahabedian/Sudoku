package text

import "testing"
import "bytes"
import "io"

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

	for _, err := range p.CheckIntegrity() {
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
	if !p.IsSolved() {
		t.Errorf("Not solved")
	}
}

func TestDiscoParty(t *testing.T) {
	// This test is from the "Disco Party" Ken-Ken published on
	// page 10 of MIT's student newspaper The Tech on 2017-09-21.
	p, err := TextToKenKen(`
		abccdd
		abccee
		affcgg
		5ffhg1
		iijhkk
		i1jllk

		a:  12 *
        b:  20 *
		c:  23 +
		d:   5 +
		e:  12 *
		f:  72 *
		g:  12 *
		h:   2 -
		i:  72 *
		j:   2 *
		k: 120 *
		l:  15 *
	`)

	if err != nil && err != io.EOF {
		t.Errorf("%s", err)
	}

	for _, err := range p.CheckIntegrity() {
		t.Errorf("%s", err)
	}

	show := func() {
		var b bytes.Buffer
		p.Show(&b)
		t.Log(b.String())
	}

	show()

	for _, g := range p.Groups {
		t.Logf("Group")
		for _, c := range g.Constraints() {
			t.Logf("constraint %s", c.Name())
		}
	}

	if err := p.DoConstraints(); err != nil {
		t.Errorf("Error during DoConstraints: %s", err.Error())
	}
	show()
	for _, j := range p.Justifications {
		t.Log(j.Pretty())
	}
	if !p.IsSolved() {
		t.Errorf("Not solved")
	}
}
