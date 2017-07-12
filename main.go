//Play your music in browser
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func main() {
	port := flag.Uint("port", 8080, "Server Port")
	path := flag.String("path", "./music", "Directotry containing your mp3 files")
	flag.Parse()
	if *port > 65535 {
		fmt.Fprintf(os.Stderr, "ERROR: Invalid port number\n")
		return
	}
	if *path == "/" {
		fmt.Fprintf(os.Stderr, "ERROR: Do not run this as root")
		return
	}
	db, filemap, err := (getSongs(*path))
	if err != nil {
		panic(err)
	}
	mux := httprouter.New()
	mux.NotFound = http.FileServer(http.Dir("public"))
	mux.HandlerFunc("GET", "/list", byTitle(db).list)
	mux.HandlerFunc("GET", "/stream", getStream(filemap))
	fmt.Printf("Serving over: http://127.0.0.1:%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))

}
