package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

type counts struct {
	*os.File
	sync.RWMutex
	db map[string]int64
}

func (count *counts) load(filename string) (err error) {
	count.db = make(map[string]int64)
	count.File, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	err = json.NewDecoder(count.File).Decode(&count.db)
	if err != nil {
		srvlog.Warn("could not decode liked json", "error", err)
	}
	return nil
}

func (count *counts) save() error {
	count.Lock()
	defer count.Unlock()
	defer count.File.Sync()
	err := count.File.Truncate(0)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(count.File)
	enc.SetIndent("", "  ")
	err = enc.Encode(count.db)
	return err
}

func (count *counts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		return
	}
	count.Lock()
	count.db[id]++
	count.Unlock()
	err := count.save()
	if err != nil {
		srvlog.Crit("Could not save likes", "error", err)
	}
}
