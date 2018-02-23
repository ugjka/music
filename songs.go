package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"sort"
	"sync"
	"time"
)

//Song object
type song struct {
	Title  string
	Artist string
	Album  string
	Track  int
	ID     string
	path   string
}

type songs []song

type idcache struct {
	db map[string]*song
	sync.RWMutex
}
type library struct {
	songs   songs
	idcache *idcache
	likes   *like
	sync.RWMutex
}

//Types for sorting
type byTitle []song
type byArtist []song
type byAlbum []song

//
// Satisfy sort interfaces
//

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

// Play counting endpoint
func countPlay(playcount map[string]int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		apiMutex.Lock()
		playcount[id]++
		savePlayCount(playcount)
		apiMutex.Unlock()
	}
}

func (cache *idcache) getCachedPaths(filemap map[string]string) {
	for id, path := range filemap {
		if _, ok := cache.db[id]; ok {
			cache.db[id].path = path
		}
	}
}

// Serves Audio files
func (cache *idcache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if v, ok := cache.db[id]; ok {
		_, err := os.Stat(v.path)
		if err != nil {
			srvlog.Crit("file missing", "file", v)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, v.path)
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (s songs) send(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	enc.Encode(s)
}

func (s songs) shuffle() {
	rand.Seed(time.Now().UnixNano())
	for i := len(s) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}
func (l *library) buildIDcache() {
	l.idcache = new(idcache)
	l.idcache.db = make(map[string]*song)
	for _, v := range l.songs {
		tmp := v
		l.idcache.db[v.ID] = &tmp
	}
}

func (l *library) getLiked() songs {
	liked := make(songs, 0)
	l.likes.RLock()
	for k := range l.likes.db {
		if v, ok := l.idcache.db[k]; ok {
			liked = append(liked, *v)
		}
	}
	l.likes.RUnlock()
	return liked
}

// Serves playlists
func (l *library) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.RLock()
	defer l.RUnlock()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	songs := l.songs
	switch r.URL.Query().Get("sort") {
	case "bytitle":
		sort.Sort(byTitle(songs))
		songs.send(w, r)
	case "byartist":
		sort.Sort(byArtist(songs))
		songs.send(w, r)
	case "byalbum":
		sort.Sort(byAlbum(songs))
		songs.send(w, r)
	case "byshuffle":
		songs.shuffle()
		songs.send(w, r)
	case "byleast":
		return
	case "bymost":
		return
	case "byfavorite":
		songs = l.getLiked()
		sort.Sort(byTitle(songs))
		songs.send(w, r)
	case "byfavshuffle":
		songs = l.getLiked()
		songs.shuffle()
		songs.send(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// Serves song/album artwork
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

// Save playcount to a file
func savePlayCount(playcount map[string]int64) {
	playcountFile.Truncate(0)
	enc := json.NewEncoder(playcountFile)
	enc.SetIndent("", " ")
	err := enc.Encode(playcount)
	if err != nil {
		srvlog.Crit("could not encode playcount json", "error", err)
	}
}
