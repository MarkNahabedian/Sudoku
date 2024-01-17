# Sudoku

Golang library for modeling and solving sudoku and kenken puzzles.

## Base

### ValueSet

We track the possible values that can appear (have not yet been eliminated) in a cell using a `ValueSet`.  Each cell has its own `ValueSet`.

`ValueSet` supports the `Len` and `IsEmpty` functions.

The `HasValue` method takes an integer as argument and returns true if
that integer is in the `ValueSet`.

The `SetHasValue` method returns a new `ValueSet` with a specified 

The `Union`, "Intersection`, and `SetDifference` methods take a second
`ValueSet` and return a new `ValueSet` that is the set union,
intersection of difference (respectively) of the two `Valueset`s.

The `DoValues` method calls a function for each value (an integer)
present in the `ValueSet`.  If the function returns false then
iteration terminates immediately without considering any further
elements of the set.


`Universe` constructs a `ValueSet` containing all possible values for
a given puzzle.  `Universe(9)` returns a `ValueSet` containing the
integers from 1 through 9 inclusive.


## Representing Puzzles

A given Sudoku or KenKen puzzle is represented by a `Puzzle` object.
When created, a puzzle has no cells or structure.
`puzzle.MakeCells(n)` adds an n by n grid of `Cell`s to `puzzle`.
Each `Cell!` has an `X` and a `Y` coordinate.  Each `Cell` also has a
list of `Group`s in which the `Cell` is a member.  Each `Cell` also
has a `ValueSet` containing the possible values that can appear in
that cell.  Values are excluded from this set as constraints are
applies.

`puzzle.AddLineGroups()` groups the cells in each row and column and
adds constraints so that no value can appear in more than one cell of
a group.

`puzzle.Add3x3Groups()` does the same, but for the 3 by 3 boxes of a
conventional Sudoku.

`Group` represents a set of cells to which some constraint
collectively applies, for example a row in a Sudoku whose cells can
not contain matching values.  Each `Group` has a list of its menber
cells and a list of constraints that apply to the cells of that group.

`puzzle.DoConstraints()` propagates constraints until exhaustion.
This will hopefully yield a solution.


## Constraints

`Constraint` is the Go interface for constraints.  Each `Constraint`
has a name and a `DoConstraint` method.

`FunctionConstraint` provides a concrete implementation of that
interface, given a name and a constraint function that is called by
its `DoConstraint` method.


### Given

"Given" is a `FunctionConstraint` that is used as the justification
when setting the single value of a `Cell` when initializing a puzzle.


### HereThenNotElsewhereConstraint,

"HereThenNotElsewhereConstraint" implements the constraint that says
that the value of one Cell of a Group can not appear in any other Cell
of that Group.


### NotElsewhereThenHereConstraint

"NotElsewhereThenHereConstraint" implements the constraint that says
that each value must appear in some Cell of the Group.


### KenKenCageConstraints

`KenKenCageConstraint` provides a concrete implementation for the
"cages" of a KenKen puzzle.

`MakeKenKenConstraint` creates a `KenKenCageConstraint` gien a
`KenKenOperator` and a value.  The constraint is satisfied if that
value equals the result of applying the operator to the value of each
cell.

These `KenKenOperator`s are supported:

#### Addition (+)

#### Subtraction (-)

#### Multiplication (*)

#### Division (/)



## Text Based Input

There is limited support for # comment lines within the grid portion
of text input.  Comments are not currently supported in the cage
definition portion of a KenLKen.

`TextToSudoku` returns an unsolved puzzle representing the specified
Sudoku.  The string argument should be a string of digits representing
the given values and dashes representing empty cells.  Spaces and tabs
are ignored.  Newlines represent breaks between rows.


`TextToKenKen` makes a ken-ken puzzle from a text specification.  The
specification starts with a grid of letters, digits and hyphens.
Cells identified by the same letter are in the same ken-ken cage.
Cells marked with a digit contain that fixed value.  Cells marked with
a hyphen are not in any cage.  After the grid description are the
rules for each cage identifying the operator and resulting value.

```
puzzle, err := TextToKenKen(`
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
```


## Command Line Solver

The text_application directory contains the source code for an
application which will solve a puzzle expressed in a text file. 
text_application/examples contains example input files.
 

## Web Based Solver

