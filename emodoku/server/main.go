// This program implements a web service which provides a browser user
// interface so that the user can design sudoku puzzles where the digits
// have been replaced by emoji symbols to produce a puzzle that can also
// artistically convey a symbolic message.
package main

import "bytes"
import "flag"
import "fmt"
import "github.com/gorilla/websocket"
import "io/ioutil"
import "log"
import "net/http"
import "os"
import "path/filepath"
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
		err = puzzle.DoConstraints()
		if err != nil {
			log.Printf("Error solving puzzle: %s", err)
			continue
		}
		b := bytes.NewBufferString("")
		puzzle.Show(b)
		log.Printf("Solved\n%s", b.String())
		
	}
}

