//Play your music in browser
package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"

	"golang.org/x/crypto/sha3"

	"github.com/julienschmidt/httprouter"
	log "gopkg.in/inconshreveable/log15.v2"
)

//Globals
var srvlog = log.New("module", "app/server")
var playcountFile *os.File
var likedFile *os.File
var sortcache = make(map[string][]byte)
var playcount = make(map[string]int64)
var liked = make(map[string]bool)
var idcache = make(map[string]bool)
var apiMutex sync.Mutex
var likedCount int
var enableFlac *bool
var passArg passwordFlag

func main() {
	//Flags
	passArg.password = base64.URLEncoding.EncodeToString(sha3.New512().Sum([]byte("")))
	flag.Var(&passArg, "password", "set password")
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
	//Create playcount store
	playcountFile, err = os.OpenFile("playcount.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(playcountFile).Decode(&playcount)
	if err != nil {
		srvlog.Warn("could not decode playcount json", "error", err)
	}
	//Create liked store
	likedFile, err = os.OpenFile("likes.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(likedFile).Decode(&liked)
	if err != nil {
		srvlog.Warn("could not decode liked json", "error", err)
	}
	//Directory for song/album art images
	os.Mkdir("artcache", 0755)
	//Parse songs
	songs, filemap := (getSongs(*path))
	//Id cache for checking stuff
	for _, v := range songs {
		idcache[v.ID] = true
	}
	//Count how many liked tracks we currently have loaded in memory
	for k := range idcache {
		if _, ok := liked[k]; ok {
			likedCount++
		}
	}
	mux := httprouter.New()
	mux.NotFound = http.FileServer(http.Dir("public"))
	mux.HandlerFunc("GET", "/count", countPlay(playcount))
	mux.HandlerFunc("GET", "/stream", getStream(filemap))
	mux.HandlerFunc("GET", "/art", artwork)
	mux.HandlerFunc("GET", "/api", getAPI(songs))
	mux.HandlerFunc("GET", "/like", likes(liked))
	mux.HandlerFunc("POST", "/login", login(passArg.password))
	fmt.Printf("Serving over: http://127.0.0.1:%d\n", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)
	srvlog.Crit("server crashed", "error", err)

}
