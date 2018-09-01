// Package html render's a base.Puzzle as an HTML table.
package html

import "bytes"
import "sudoku/base"
import "strings"
import "html/template"


func ToTable(puzzle *base.Puzzle, glyphs map[int]rune) string {
	s := &spec{ Puzzle: puzzle, Glyphs: glyphs}
	writer := bytes.NewBufferString("")
	err := table_template.Execute(writer, s)
	if err != nil {
		panic(err)
	}
	return writer.String()
}

// spec is the input to table_template.
type spec struct {
	Puzzle *base.Puzzle
	Glyphs map[int]rune
}

func (s *spec) RowIndices() (indices []int) {
	for i := 1; i <= s.Puzzle.Size; i++ {
		indices = append(indices, i)
	}
	return indices
}

func (s *spec) ColumnIndices() (indices []int) {
	for i := 1; i <= s.Puzzle.Size; i++ {
		indices = append(indices, i)
	}
	return indices
}

func (s *spec) Glyph(rowIndex, columnIndex int) string {
	cell := s.Puzzle.Cell(columnIndex, rowIndex)
	isSolved, value := cell.IsSolved()
	if !isSolved {
		return ""
	}
	glyph, ok := s.Glyphs[value]
	if ok {
		return string([]rune{glyph})
	}
	return "0123456789"[value:value+1]
}

// BorderClass returns the value for the HTML CSS class attribute
// for a TD element to specify how to draw its borders.
func (s *spec) BorderClass(rowIndex, columnIndex int) string {
	classes := []string{}
	ri := rowIndex - 1
	ci := columnIndex - 1
	switch ri % 3 {
	case 0:
		classes = append(classes, "top")
	case 1:
		classes = append(classes, "vmiddle")
	case 2:
		classes = append(classes, "bottom")
	}
	switch ci % 3 {
	case 0:
		classes = append(classes, "left")
	case 1:
		classes = append(classes, "hmiddle")
	case 2:
		classes = append(classes, "right")
	}
	return strings.Join(classes, " ")
}


var table_template = template.Must(template.New("name").Parse(`
{{with $spec := .}}
	<table>
		{{range $rowIndex := $spec.RowIndices}}
			<tr>
				{{range $columnIndex := $spec.ColumnIndices}}
					<td id="row{{$rowIndex}}_col{{$columnIndex}}"
					    class="{{$spec.BorderClass $rowIndex $columnIndex}}">
						{{$spec.Glyph $rowIndex $columnIndex}}
					</td>
				{{end}}
			</tr>
		{{end}}
	</table>
{{end}}
`))
