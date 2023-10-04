package router

import "net/http"

// Shortcut methods for common HTTP methods

func (r *Router) GET(path string, handler http.HandlerFunc) {
    r.Handle("GET", path, handler)
}

func (r *Router) POST(path string, handler http.HandlerFunc) {
    r.Handle("POST", path, handler)
}

func (r *Router) PUT(path string, handler http.HandlerFunc) {
    r.Handle("PUT", path, handler)
}

func (r *Router) DELETE(path string, handler http.HandlerFunc) {
    r.Handle("DELETE", path, handler)
}

// Add similar functions for other HTTP methods as needed
