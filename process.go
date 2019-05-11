package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

//
//Process audio files, extract info/artwork
//

func getSongs(searchdir string, flac bool) (songs songs) {
	cache := make(map[string]song)
	cachef, err := os.OpenFile("cache.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer cachef.Close()
	err = json.NewDecoder(cachef).Decode(&cache)
	if err != nil {
		srvlog.Warn("could not decode cache json", "error", err)
	}
	filepath.Walk(searchdir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !(strings.HasSuffix(path, ".mp3") || (strings.HasSuffix(path, ".flac") && flac)) {
			//srvlog.Info("skipping invalid file", "file", info.Name())
			return nil
		}
		if v, ok := cache[path]; ok {
			v.path = path
			songs = append(songs, v)
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
			srvlog.Warn("could not read tags", "file", info.Name())
			return nil
		}
		artist := strings.TrimSpace(m.Artist())
		title := strings.TrimSpace(m.Title())
		album := strings.TrimSpace(m.Album())
		track, _ := m.Track()
		hash, err := tag.Sum(f)
		if err != nil {
			srvlog.Warn("could not hash the file", "file", info.Name())
			return nil
		}
		if m.Picture() != nil {
			art, err := os.OpenFile(fmt.Sprintf("./artcache/%s", hash), os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0644)
			if err != nil {
				panic(err)
			}
			io.Copy(art, bytes.NewReader(m.Picture().Data))
			art.Close()
		}
		result := song{Artist: artist, Title: title, Album: album, Track: track, ID: hash, path: path}
		songs = append(songs, result)
		cache[path] = result
		return nil
	})
	err = cachef.Truncate(0)
	if err != nil {
		srvlog.Warn("could not truncate cache file", "error", err)
		return
	}
	enc := json.NewEncoder(cachef)
	enc.SetIndent("", "\t")
	err = enc.Encode(cache)
	if err != nil {
		srvlog.Warn("could not encode cache json", "error", err)
	}
	return
}
