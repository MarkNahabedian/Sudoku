// Define some metrics to measure how colse a Puzzle is to being
// solved.
package base

// Returns the number of possible values summed over all cells.
func (p *Puzzle) ValueCount() int {
	count := 0
	for _, cell := range p.Grid {
		count += cell.Possibilities.Len()
	}
	return count
}

// Returns what the number of possible values summed over all cells
// would be when the puzzle is solved: that would be one value per
// cell.
func (p *Puzzle) SolvedValueCount() int {
	return p.Size * p.Size
}

// Returns what the number of possible values summed over all cells
// would be if we knew nothing about the puzzle: number of cells times
// the size of the universe.
func (p *Puzzle) MaxValueCount() int {
	return p.Size * p.Size * p.Universe.Len()
}

