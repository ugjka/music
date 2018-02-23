package main

import (
	"net/http"
	"sort"
	"sync"
)

type library struct {
	songs   songs
	idcache *idcache
	likes   *like
	counts  *counts
	sync.RWMutex
}

func (l *library) makeCache() {
	l.idcache = new(idcache)
	l.idcache.db = make(map[string]*song)
	for _, v := range l.songs {
		tmp := v
		l.idcache.db[v.ID] = &tmp
	}
}

func (l *library) Liked() songs {
	liked := make(songs, 0)
	l.likes.RLock()
	l.idcache.RLock()
	for k := range l.likes.db {
		if v, ok := l.idcache.db[k]; ok {
			liked = append(liked, *v)
		}
	}
	l.idcache.RUnlock()
	l.likes.RUnlock()
	return liked
}

func (l *library) Counted() songs {
	counted := make(songs, 0)
	l.counts.RLock()
	l.idcache.RLock()
	for id, song := range l.idcache.db {
		tmp := *song
		if count, ok := l.counts.db[id]; ok {
			tmp.playcount = count
		}
		counted = append(counted, tmp)
	}
	l.idcache.RUnlock()
	l.counts.RUnlock()
	return counted
}

// Serves playlists
func (l *library) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	songs := make(songs, len(l.songs))
	l.RLock()
	copy(songs, l.songs)
	l.RUnlock()
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
		songs = l.Counted()
		sort.Sort(byPlayCount(songs))
		songs.send(w, r)
	case "bymost":
		songs = l.Counted()
		sort.Sort(sort.Reverse(byPlayCount(songs)))
		songs.send(w, r)
	case "byfavorite":
		songs = l.Liked()
		sort.Sort(byTitle(songs))
		songs.send(w, r)
	case "byfavshuffle":
		songs = l.Liked()
		songs.shuffle()
		songs.send(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
