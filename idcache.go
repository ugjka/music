package main

import (
	"net/http"
	"os"
	"sync"
)

type idcache struct {
	db map[string]*song
	sync.RWMutex
}

func (cache *idcache) getCachedPaths(filemap map[string]string) {
	for id, path := range filemap {
		if _, ok := cache.db[id]; ok {
			cache.db[id].path = path
		}
	}
}

// Serves Audio files
func (cache *idcache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	cache.RLock()
	v, ok := cache.db[id]
	cache.RUnlock()
	if ok {
		_, err := os.Stat(v.path)
		if err != nil {
			srvlog.Crit("file missing", "file", v)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, v.path)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	return

}
