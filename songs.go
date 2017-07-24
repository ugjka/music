package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"sort"
	"sync"
	"time"
)

type song struct {
	Title  string
	Artist string
	Album  string
	Track  int
	ID     string
	path   string
}

type byTitle []*song
type byArtist []*song
type byAlbum []*song

func (songs byTitle) Len() int {
	return len(songs)
}

func (songs byArtist) Len() int {
	return len(songs)
}

func (songs byAlbum) Len() int {
	return len(songs)
}

func (songs byTitle) Swap(i, j int) {
	songs[i], songs[j] = songs[j], songs[i]
}

func (songs byArtist) Swap(i, j int) {
	songs[i], songs[j] = songs[j], songs[i]
}

func (songs byAlbum) Swap(i, j int) {
	songs[i], songs[j] = songs[j], songs[i]
}

func (songs byTitle) Less(i, j int) bool {
	if songs[i].Title != songs[j].Title {
		return songs[i].Title < songs[j].Title
	}
	if songs[i].Artist != songs[j].Artist {
		return songs[i].Artist < songs[j].Artist
	}
	if songs[i].Album != songs[j].Album {
		return songs[i].Album < songs[j].Album
	}
	if songs[i].Track != songs[j].Track {
		return songs[i].Track < songs[j].Track
	}
	return false
}

func (songs byArtist) Less(i, j int) bool {
	if songs[i].Artist != songs[j].Artist {
		return songs[i].Artist < songs[j].Artist
	}
	if songs[i].Album != songs[j].Album {
		return songs[i].Album < songs[j].Album
	}
	if songs[i].Track != songs[j].Track {
		return songs[i].Track < songs[j].Track
	}
	if songs[i].Title != songs[j].Title {
		return songs[i].Title < songs[j].Title
	}
	return false
}

func (songs byAlbum) Less(i, j int) bool {
	if songs[i].Album != songs[j].Album {
		return songs[i].Album < songs[j].Album
	}
	if songs[i].Artist != songs[j].Artist {
		return songs[i].Artist < songs[j].Artist
	}
	if songs[i].Track != songs[j].Track {
		return songs[i].Track < songs[j].Track
	}
	if songs[i].Title != songs[j].Title {
		return songs[i].Title < songs[j].Title
	}
	return false
}

func getStream(filemap map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if v, ok := filemap[r.URL.Query().Get("id")]; ok {
			_, err := os.Stat(v)
			if err != nil {
				srvlog.Crit("file missing", "file", v)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			http.ServeFile(w, r, v)
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

var sortcache = make(map[string][]byte)

var apiMutex sync.Mutex

func getAPI(songs []*song) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var err error
		switch r.URL.Query().Get("sort") {
		case "bytitle":
			if v, ok := sortcache["bytitle"]; !ok {
				apiMutex.Lock()
				sort.Sort(byTitle(songs))
				enc := json.NewEncoder(w)
				enc.SetIndent("", " ")
				enc.Encode(songs)
				sortcache["bytitle"], err = json.MarshalIndent(songs, "", " ")
				apiMutex.Unlock()
				if err != nil {
					srvlog.Crit("marshaling cache by title failed", "error", err)
					delete(sortcache, "bytitle")
				}
			} else {
				io.Copy(w, bytes.NewReader(v))
			}
		case "byartist":
			if v, ok := sortcache["byartist"]; !ok {
				apiMutex.Lock()
				sort.Sort(byArtist(songs))
				enc := json.NewEncoder(w)
				enc.SetIndent("", " ")
				enc.Encode(songs)
				sortcache["byartist"], err = json.MarshalIndent(songs, "", " ")
				apiMutex.Unlock()
				if err != nil {
					srvlog.Crit("marshaling cache by artist failed", "error", err)
					delete(sortcache, "byartist")
				}
			} else {
				io.Copy(w, bytes.NewReader(v))
			}
		case "byalbum":
			if v, ok := sortcache["byalbum"]; !ok {
				apiMutex.Lock()
				sort.Sort(byAlbum(songs))
				enc := json.NewEncoder(w)
				enc.SetIndent("", " ")
				enc.Encode(songs)
				sortcache["byalbum"], err = json.MarshalIndent(songs, "", " ")
				apiMutex.Unlock()
				if err != nil {
					srvlog.Crit("marshaling cache by album failed", "error", err)
					delete(sortcache, "byalbum")
				}
			} else {
				io.Copy(w, bytes.NewReader(v))
			}
		case "byshuffle":
			apiMutex.Lock()
			rand.Seed(time.Now().UnixNano())
			for i := len(songs) - 1; i > 0; i-- {
				j := rand.Intn(i + 1)
				songs[i], songs[j] = songs[j], songs[i]
			}
			enc := json.NewEncoder(w)
			enc.SetIndent("", " ")
			enc.Encode(songs)
			apiMutex.Unlock()
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func artwork(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	reg := regexp.MustCompile("^\\w+$")
	if !reg.MatchString(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	file := fmt.Sprintf("./artcache/%s", id)
	_, err := os.Stat(file)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, file)
}
