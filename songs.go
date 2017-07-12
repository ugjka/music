package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"
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

func (db byTitle) list(w http.ResponseWriter, req *http.Request) {
	tw := new(tabwriter.Writer).Init(w, 0, 8, 2, ' ', 0)
	for _, v := range db {
		fmt.Fprintf(tw, "%v\t%v\t%v\t%v\t%v\t\n", v.Artist, v.Title, v.Album, v.Track, v.ID)
	}
	tw.Flush()
}

func getStream(filemap map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if v, ok := filemap[r.URL.Query().Get("hash")]; ok {
			f, err := os.Open(v)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			defer f.Close()
			w.Header().Set("Content-Type", "audio/mpeg")
			io.Copy(w, f)
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}
