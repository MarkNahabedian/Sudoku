// Package main implements a text based application for solving sudoku
// and kenken puzzles.
package main

import "sudoku/base"
import "sudoku/text"
import "flag"
import "fmt"
import "io/ioutil"
import "os"
import "strings"

type PuzzleType struct {
	Name string
	Parser func(text string) (*base.Puzzle, error)
	Example string
}

func find_puzzle_type(name string) (*PuzzleType, error) {
	supported := []string{}
	for _, pt := range PuzzleTypes {
		supported = append(supported, pt.Name)
		if pt.Name == name {
			return pt, nil
		}
	}
	return nil, fmt.Errorf("The only supported values for the --puzzle flag are %s", strings.Join(supported, ", "))
}

func (pt *PuzzleType) String() string {
	return pt.Name
}

type PuzzleTypeVar struct {
	Value *PuzzleType
}

func (ptv *PuzzleTypeVar) String() string {
	if ptv.Value == nil {
		return ""
	}
	return ptv.Value.Name
}

func (ptv *PuzzleTypeVar) Set(name string) error {
	pt, err := find_puzzle_type(name)
	if err != nil {
		return err
	}
	ptv.Value = pt
	return nil
}

func (ptv *PuzzleTypeVar) Get() *PuzzleType {
	return ptv.Value
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(-1)
}

func main() {
	var input string
	var output string
	var puzzle_type PuzzleTypeVar

	flag.StringVar(&input, "input", "", "Path to a file containing the unsolved puzzle.")
	flag.StringVar(&output, "output", "-", "The file to write the solved puzzle to.")
	flag.Var(&puzzle_type, "puzzle",
		"The type of puzzle to solve, either 'sudoku' or 'kenken'.  If not specified an example puzzle is used.")
	flag.Parse()

	/*
	fmt.Printf("flags.Visit:\n")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("--%s=%v  %s\n", f.Name, f.Value, f.Usage)
	})
	fmt.Printf("Positional arguments: %v\n", flag.Args())
	*/

	if puzzle_type.Value == nil {
		fmt.Printf("puzzle_type not set.  Available puzzle types: ")
		for i, pt := range PuzzleTypes {
			if i > 0 { fmt.Printf(", ") }
			fmt.Printf("%s", pt.Name)
		}
		fmt.Printf("\n")
		os.Exit(-1)
	}

	var puzzle_string string
	if input == "" {
		puzzle_string = puzzle_type.Value.Example
	} else {
		bytes, err := ioutil.ReadFile(input)
		if err != nil {
			fail(fmt.Errorf("Can't read %s: %s", input, err))
		}
		puzzle_string = string(bytes)
	}

	puzzle, err := puzzle_type.Value.Parser(puzzle_string)
	if err != nil {
		border := strings.Repeat("=", 30)
		fmt.Printf("%s\n%s\n%s\n",
			border, puzzle_string, border)
		fail(err)
	}

	var out *os.File
	if output == "-" {
		out = os.Stdout
	} else {
		out, err := os.Create(output)
		if err != nil {
			fail(fmt.Errorf("Can't open %s: %s", output, err))
		}
		defer out.Close()
	}

	// First write the original unsolved puzzle.
	out.WriteString(puzzle_string)

	pre_solve_value_count := puzzle.ValueCount()

	// Solve it
	err = puzzle.GuessSolve()

	// Write the answer
	puzzle.Show(out)

	if !puzzle.IsSolved() {
		fmt.Printf("Progress: %d %d %d %d\n\n",
			puzzle.MaxValueCount(),
			pre_solve_value_count,
			puzzle.ValueCount(),
			puzzle.SolvedValueCount())
	}
	
	// Write the justifications.
	for _, j := range puzzle.Justifications {
		out.WriteString(j.Pretty())
		out.WriteString("\n")
	}

	if err != nil {
		fail(fmt.Errorf("Error while solving: %s", err.Error()))
	}
}


var PuzzleTypes = []*PuzzleType {
	&PuzzleType{
		Name: "sudoku",
		Parser: text.TextToSudoku,
		Example: `
	---7-----
	1--------
	---43-2--
	--------6
	---5-9---
	------418
	----81---
	--2----5-
	-4----3--
	`},
	&PuzzleType{
		Name: "kenken",
		Parser: text.TextToKenKen,
		Example: `
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
	`},
}
