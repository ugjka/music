package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
)

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
