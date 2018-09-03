// This program implements a web service which provides a browser user
// interface so that the user can design sudoku puzzles where the digits
// have been replaced by emoji symbols to produce a puzzle that can also
// artistically convey a symbolic message.
package main

import "encoding/json"
import "flag"
import "fmt"
import "github.com/gorilla/websocket"
import "io/ioutil"
import "log"
import "net/http"
import "os"
import "path/filepath"
import "sudoku/base"
import "sudoku/text"

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	// ***** Should replace this with something safer some day.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var webAddress = flag.String("server-address", ":8000",
	"The address for this web server.")

var topStaticFile = flag.String("top", "emodoku.html",
	"The path to the static top level page.")

var bordersCSS = flag.String("bordersCSS", "../../html/borders.css",
	"The path to the borders.css file.")

var jsFile = flag.String("jsFile", "emodoku.js",
	"The path to the client-side javascript file.")

var makeSudoku = flag.String("makeSudoku", "make_sudoku.html",
	"The path to the html page that implements the service.")

func main() {
	flag.Parse()
	http.HandleFunc("/", makeFileResponder(*topStaticFile))
	http.HandleFunc("/borders.css", makeFileResponder(*bordersCSS))
	http.HandleFunc("/emodoku.js", makeFileResponder(*jsFile))
	http.HandleFunc("/make-sudoku.html", makeFileResponder(*makeSudoku))
	http.HandleFunc("/solver", handleSolver)
	err := http.ListenAndServe(*webAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func makeFileResponder(file string) func(http.ResponseWriter, *http.Request) {
	// http.ServeFile doesn't accept .. for security reasons.
	file, err := filepath.Abs(file)
	// Don't start the server if the file doesn't exist.
	_, err = os.Stat(file)
	if err != nil {
		panic(fmt.Sprintf("File %s: %s", file, err))
	}
	if err != nil {
		panic(err)
	}
	return func (w http.ResponseWriter, r *http.Request) {
		log.Printf("ServeFile %s", file)
		http.ServeFile(w, r, file)
	}
}

func handleSolver(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %s", err)
		return
	}
	for {
		messageType, r, err := conn.NextReader()
		if err != nil {
			log.Printf("conn.NextReader: %s", err)
			return
		}
		if messageType != websocket.TextMessage  {
			log.Printf("Received unsupported message type %d", messageType)
			continue
		}
		msg, err := ioutil.ReadAll(r)
		if err != nil {
			log.Printf("Error reading message: %s", err)
			continue
		}
		puzzle, err := text.TextToSudoku(string(msg))
		if err != nil {
			log.Printf("Error parsing sudoku from text: %s", err)
			continue
		}
		response := MakeSolutionResponse(puzzle)
		encoded, err := json.MarshalIndent(response, "", "")
		if err != nil {
			log.Printf("Error ncoding JSON: %s", err)
			continue
		}
		log.Printf("Solver response:\n%s", encoded)
		w, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Printf("NextWriter error: %s",err)
			continue
		}
		_, err = w.Write(encoded)
		if err != nil {
			log.Printf("Error writing to socket: %s", err)
		}
	}
}

type SolutionResponse struct {
	// Size is the number of rows and columns.
	Size uint
	// Possibilities has an array of ints for eeach cell of the puzzle grid.
	// This first index is the row number, the second the column number.
	// these indices are zero origin because that's how vectors work in the
	// languages we're using.
	Possibilities [][][]uint
	// Error will be the empty string if no error occurred while the puzzle
	// was being solved, otherwise it is a string describing the error.
	Error string
}

func MakeSolutionResponse(p *base.Puzzle) *SolutionResponse {
	errMsg := ""
	err := p.DoConstraints()
	if err != nil {
		errMsg = err.Error()
	}
	grid := make([][][]uint, p.Size)
	for row := 0; row < p.Size; row++ {
		grid[row] = make([][]uint, p.Size)
		for col := 0; col < p.Size; col++ {
			cell := p.Cell(col + 1, row + 1)
			for val := 1; val <= p.Size; val++ {
				if cell.Possibilities.HasValue(val) {
					grid[row][col] = append(grid[row][col], uint(val))
				}
			}
		}
	}
	return &SolutionResponse{
		Size: uint(p.Size),
		Possibilities: grid,
		Error: errMsg,
	}
}
