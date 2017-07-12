package main

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
