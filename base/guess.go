// If constraint propagation doesn't solve the puzzle, try iterative
// guessing.
package base

var Pick = FunctionConstraint{
	name:       "Pick",
	constraint: (func(g *Group) error { return nil }),
}

// Find the first Cell that isn't solved and pick a value for it.
func (p *Puzzle) GuessOnce() error {
	if p.IsSolved() {
		return nil
	}
	var guess_err error = nil
	for _, cell := range p.Grid {
		if s, _ := cell.IsSolved(); !s {
			cell.Possibilities.DoValues(
				func(value int) bool {
					_, guess_err = cell.MustBe(value, Pick, nil)
					return false
				})
			if guess_err != nil {
				return guess_err
			}
			break
		}
	}
	return nil
}

// Try to solve the puzzle by guessing.
func (p *Puzzle) GuessSolve() error {
	for {
		if err := p.DoConstraints(); err != nil {
			return err
		}
		if p.IsSolved() {
			break
		}
		if err := p.GuessOnce(); err != nil {
			return err
		}
	}
	return nil
}
