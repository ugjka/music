//Play your music in browser
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	log "gopkg.in/inconshreveable/log15.v2"
)

var srvlog = log.New("module", "app/server")

func main() {
	port := flag.Uint("port", 8080, "Server Port")
	path := flag.String("path", "./music", "Directory containing your mp3 files")
	flag.Parse()
	if *port > 65535 {
		fmt.Fprintf(os.Stderr, "ERROR: Invalid port number\n")
		return
	}
	if *path == "/" {
		fmt.Fprintf(os.Stderr, "ERROR: Do not run this as root")
		return
	}
	songs, filemap := (getSongs(*path))
	mux := httprouter.New()
	mux.NotFound = http.FileServer(http.Dir("public"))
	mux.HandlerFunc("GET", "/stream", getStream(filemap))
	mux.HandlerFunc("GET", "/api", getAPI(songs))
	fmt.Printf("Serving over: http://127.0.0.1:%d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)
	srvlog.Crit("server crashed", "error", err)

}
