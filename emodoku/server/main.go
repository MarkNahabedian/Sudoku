// This program implements a web service which provides a browser user
// interface so that the user can design sudoku puzzles where the digits
// have been replaced by emoji symbols to produce a puzzle that can also
// artistically convey a symbolic message.
package main

import "path/filepath"
import "flag"
import "fmt"
import "log"
import "net/http"
import "os"
import "github.com/gorilla/websocket"

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
	// This will eventually call a handler function that the client will do websocket interactions with.
	http.HandleFunc("/make-sudoku.html", makeFileResponder(*makeSudoku))
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


