package router

import (
    "net/http"
)

type Router struct {
    // holds the routes
    routes []Route
}

type Route struct {
    path    string
    method  string
    handler http.HandlerFunc
}

func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
    route := Route{
        path:    path,
        method:  method,
        handler: handler,
    }
    r.routes = append(r.routes, route)
}
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    for _, route := range r.routes {
        if route.path == req.URL.Path && route.method == req.Method {
            route.handler(w, req)
            return
        }
    }
    http.NotFound(w, req)
}
func NewRouter() *Router {
    return &Router{}
}
