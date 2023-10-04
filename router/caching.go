package router

import (
	"bytes"
	"net/http"
	"sync"
	"time"
)

type cacheItem struct {
	body       []byte
	header     http.Header
	expiration time.Time
}

var (
	cache = make(map[string]*cacheItem)
	lock  = sync.RWMutex{}
)

func CachingMiddleware(duration time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" { // Only cache GET requests
			next.ServeHTTP(w, r)
			return
		}

		lock.RLock()
		item, exists := cache[r.RequestURI]
		lock.RUnlock()

		if exists && time.Now().Before(item.expiration) {
			for k, v := range item.header {
				w.Header()[k] = v
			}
			w.Write(item.body)
			return
		}

		cw := &cacheWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			header:         http.Header{},
		}

		next.ServeHTTP(cw, r)

		if cw.statusCode == http.StatusOK { // Only cache successful responses
			lock.Lock()
			cache[r.RequestURI] = &cacheItem{
				body:       cw.body.Bytes(),
				header:     cw.header,
				expiration: time.Now().Add(duration),
			}
			lock.Unlock()
		}
	})
}

type cacheWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
	header     http.Header
}

func (cw *cacheWriter) WriteHeader(statusCode int) {
	cw.statusCode = statusCode
	cw.ResponseWriter.WriteHeader(statusCode)
}

func (cw *cacheWriter) Write(b []byte) (int, error) {
	if cw.body == nil {
		cw.body = bytes.NewBuffer(b)
	} else {
		cw.body.Write(b)
	}
	return cw.ResponseWriter.Write(b)
}

func (cw *cacheWriter) Header() http.Header {
	return cw.header
}
