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

//Types for sorting
type byTitle []*song
type byArtist []*song
type byAlbum []*song
type byLeast []*song
type byFavorite []*song

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

func (songs byLeast) Len() int {
	return len(songs)
}

func (songs byFavorite) Len() int {
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

func (songs byLeast) Swap(i, j int) {
	songs[i], songs[j] = songs[j], songs[i]
}

func (songs byFavorite) Swap(i, j int) {
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

func (songs byLeast) Less(i, j int) bool {
	var icount, jcount int64
	if v, ok := playcount[songs[i].ID]; ok {
		icount = v
	} else {
		icount = 0
	}
	if v, ok := playcount[songs[j].ID]; ok {
		jcount = v
	} else {
		jcount = 0
	}
	if icount != jcount {
		return icount < jcount
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
	if songs[i].Title != songs[j].Title {
		return songs[i].Title < songs[j].Title
	}
	return false
}

func (songs byFavorite) Less(i, j int) bool {
	var icount, jcount int
	if _, ok := liked[songs[i].ID]; ok {
		icount = 1
	} else {
		icount = 0
	}
	if _, ok := liked[songs[j].ID]; ok {
		jcount = 1
	} else {
		jcount = 0
	}
	if icount != jcount {
		return icount > jcount
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

// Get likes or set or remove likes
func likes(likes map[string]bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiMutex.Lock()
		defer apiMutex.Unlock()
		if id := r.URL.Query().Get("id"); id != "" {
			if v, ok := likes[id]; ok {
				json.NewEncoder(w).Encode(v)
			} else {
				json.NewEncoder(w).Encode(false)
			}
			return
		}
		if like := r.URL.Query().Get("like"); like != "" {
			if _, ok := idcache[like]; !ok {
				return
			}
			if _, ok := likes[like]; ok {
				delete(likes, like)
				likedCount--
				json.NewEncoder(w).Encode(false)
			} else {
				likes[like] = true
				likedCount++
				json.NewEncoder(w).Encode(true)
			}
		}
		err := likedFile.Truncate(0)
		if err != nil {
			srvlog.Crit("Could not truncate likes file", "error", err)
			return
		}
		enc := json.NewEncoder(likedFile)
		enc.SetIndent("", " ")
		err = enc.Encode(likes)
		if err != nil {
			srvlog.Crit("Could not encode likes json", "error", err)
		}
	}
}

// Serves Audio files
func getStream(filemap map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if v, ok := filemap[id]; ok {
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

// Serves playlists
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
		case "byleast":
			apiMutex.Lock()
			sort.Sort(byLeast(songs))
			enc := json.NewEncoder(w)
			enc.SetIndent("", " ")
			enc.Encode(songs)
			apiMutex.Unlock()
		case "bymost":
			apiMutex.Lock()
			sort.Sort(sort.Reverse(byLeast(songs)))
			enc := json.NewEncoder(w)
			enc.SetIndent("", " ")
			enc.Encode(songs)
			apiMutex.Unlock()
		case "byfavorite":
			apiMutex.Lock()
			sort.Sort(byFavorite(songs))
			enc := json.NewEncoder(w)
			enc.SetIndent("", " ")
			enc.Encode(songs[0:likedCount])
			apiMutex.Unlock()
		case "byfavshuffle":
			apiMutex.Lock()
			sort.Sort(byFavorite(songs))
			rand.Seed(time.Now().UnixNano())
			for i := len(songs[0:likedCount]) - 1; i > 0; i-- {
				j := rand.Intn(i + 1)
				songs[i], songs[j] = songs[j], songs[i]
			}
			enc := json.NewEncoder(w)
			enc.SetIndent("", " ")
			enc.Encode(songs[0:likedCount])
			apiMutex.Unlock()
		default:
			w.WriteHeader(http.StatusNotFound)
		}
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
