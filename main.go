//Play your music in browser
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/crypto/sha3"

	"github.com/julienschmidt/httprouter"
	log "gopkg.in/inconshreveable/log15.v2"
)

//Globals
var srvlog = log.New("module", "app/server")
var enableFlac *bool
var pass passwordFlag

func main() {
	//Flags
	pass.password = base64.URLEncoding.EncodeToString(sha3.New512().Sum([]byte("")))
	flag.Var(&pass, "password", "set password")
	enableFlac = flag.Bool("enableFlac", false, "Enable Flac streaming")
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
	//Directory for song/album art images
	os.Mkdir("artcache", 0755)
	//Parse songs
	songs, filemap := getSongs(*path)

	var likes = &like{}
	err = likes.load("./likes.json")
	if err != nil {
		panic(err)
	}

	var counts = &count{}
	err = counts.load("./playcount.json")
	if err != nil {
		panic(err)
	}

	var library = &library{}
	library.songs = songs
	library.likes = likes
	library.counts = counts
	library.makeCache()
	library.idcache.getPaths(filemap)

	mux := httprouter.New()
	mux.NotFound = http.FileServer(http.Dir("public"))
	mux.Handler("GET", "/count", pass.mustAuth(library.counts))
	mux.Handler("GET", "/stream", pass.mustAuth(library.idcache))
	mux.HandlerFunc("GET", "/art", artwork)
	mux.Handler("GET", "/api", library)
	mux.Handler("GET", "/like", pass.mustAuth(likes))
	mux.HandlerFunc("POST", "/login", login(pass.password))
	fmt.Printf("Serving over: http://127.0.0.1:%d\n", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)
	srvlog.Crit("server crashed", "error", err)

}
