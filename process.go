package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

func getSongs(searchdir string) (songs []*song, filemap map[string]string) {
	filemap = make(map[string]string)
	filepath.Walk(searchdir, func(path string, finfo os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".mp3") {
			srvlog.Info("skipping non-mp3 file", "file", finfo.Name())
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			srvlog.Warn("could not read file", "path", path)
			return nil
		}
		defer f.Close()
		m, err := tag.ReadFrom(f)
		if err != nil {
			srvlog.Warn("could not read tags", "file", finfo.Name())
			return nil
		}
		artist := strings.TrimSpace(m.Artist())
		title := strings.TrimSpace(m.Title())
		album := strings.TrimSpace(m.Album())
		track, _ := m.Track()
		hash, err := tag.Sum(f)
		if err != nil {
			srvlog.Warn("could not hash the file", "file", finfo.Name())
			return nil
		}
		filemap[hash] = path
		songs = append(songs, &song{Artist: artist, Title: title, Album: album, Track: track, ID: hash, path: path})
		return nil
	})
	return
}
