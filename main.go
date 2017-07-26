//Play your music in browser
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/julienschmidt/httprouter"
	log "gopkg.in/inconshreveable/log15.v2"
)

var srvlog = log.New("module", "app/server")
var playcountFile *os.File
var sortcache = make(map[string][]byte)
var playcount = make(map[string]int64)
var apiMutex sync.Mutex

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
	var err error
	playcountFile, err = os.OpenFile("playcount.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(playcountFile).Decode(&playcount)
	if err != nil {
		srvlog.Warn("could not decode playcount json", "error", err)
	}
	os.Mkdir("artcache", 0755)
	songs, filemap := (getSongs(*path))
	mux := httprouter.New()
	mux.NotFound = http.FileServer(http.Dir("public"))
	mux.HandlerFunc("GET", "/count", countPlay(playcount))
	mux.HandlerFunc("GET", "/stream", getStream(filemap))
	mux.HandlerFunc("GET", "/art", artwork)
	mux.HandlerFunc("GET", "/api", getAPI(songs))
	fmt.Printf("Serving over: http://127.0.0.1:%d\n", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)
	srvlog.Crit("server crashed", "error", err)

}
