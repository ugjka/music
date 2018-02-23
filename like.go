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

func (like *like) load(filename string) (err error) {
	like.db = make(map[string]bool)
	like.File, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
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
		_, ok := like.db[id]
		like.RUnlock()
		json.NewEncoder(w).Encode(ok)
		return
	}
	if liked := r.URL.Query().Get("like"); liked != "" {
		if _, ok := like.db[liked]; ok {
			like.Lock()
			delete(like.db, liked)
			like.Unlock()
			json.NewEncoder(w).Encode(false)
		} else {
			like.Lock()
			like.db[liked] = true
			like.Unlock()
			json.NewEncoder(w).Encode(true)
		}
		err := like.save()
		if err != nil {
			srvlog.Crit("Could not save likes", "error", err)
		}
	}
}
