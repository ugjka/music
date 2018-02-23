package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

type like struct {
	db map[string]bool
	*os.File
	sync.RWMutex
}

func (like *like) loadLikes(filename string) (err error) {
	like.db = make(map[string]bool)
	like.File, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	err = json.NewDecoder(like.File).Decode(&like.db)
	if err != nil {
		srvlog.Warn("could not decode liked json", "error", err)
	}
	return nil
}

func (like *like) save() error {
	like.Lock()
	defer like.Unlock()
	defer like.File.Sync()
	err := like.File.Truncate(0)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(like.File)
	enc.SetIndent("", "  ")
	err = enc.Encode(like.db)
	return err
}

// Get likes or set or remove likes
func (like *like) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if id := r.URL.Query().Get("id"); id != "" {
		like.RLock()
		defer like.RUnlock()
		_, ok := like.db[id]
		json.NewEncoder(w).Encode(ok)
		return
	}
	if liked := r.URL.Query().Get("like"); liked != "" {
		like.Lock()
		if _, ok := like.db[liked]; ok {
			delete(like.db, liked)
			json.NewEncoder(w).Encode(false)
		} else {
			like.db[liked] = true
			json.NewEncoder(w).Encode(true)
		}
		like.Unlock()
		err := like.save()
		if err != nil {
			srvlog.Crit("Could not save likes", "error", err)
		}
	}
}
