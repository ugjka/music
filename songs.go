package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
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

func getAPI(songs []*song) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		switch r.URL.Query().Get("sort") {
		case "bytitle":
			sort.Sort(byTitle(songs))
			enc := json.NewEncoder(w)
			enc.SetIndent("", " ")
			enc.Encode(songs)
		case "byartist":
			sort.Sort(byArtist(songs))
			enc := json.NewEncoder(w)
			enc.SetIndent("", " ")
			enc.Encode(songs)
		case "byalbum":
			sort.Sort(byAlbum(songs))
			enc := json.NewEncoder(w)
			enc.SetIndent("", " ")
			enc.Encode(songs)
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
