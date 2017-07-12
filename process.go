package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

func getSongs(searchdir string) (songs []*song, filemap map[string]string, err error) {
	filemap = make(map[string]string)
	err = filepath.Walk(searchdir, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".mp3") {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		m, err := tag.ReadFrom(f)
		if err != nil {
			return nil
		}
		artist := strings.TrimSpace(m.Artist())
		title := strings.TrimSpace(m.Title())
		album := strings.TrimSpace(m.Album())
		track, _ := m.Track()
		hash, err := tag.Sum(f)
		if err != nil {
			return nil
		}
		filemap[hash] = path
		songs = append(songs, &song{Artist: artist, Title: title, Album: album, Track: track, ID: hash, path: path})
		return nil
	})
	return
}
