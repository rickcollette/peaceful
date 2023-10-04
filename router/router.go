package router

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Router struct {
	routes     []Route
	middleware MiddlewareChain // Changed to MiddlewareChain
	routeGroups map[string]*RouteGroup
}

type Route struct {
	path       string
	method     string
	handler    http.HandlerFunc
	pattern    *regexp.Regexp
	params     map[int]string
	middleware MiddlewareChain // Changed to MiddlewareChain
}

type Middleware func(http.Handler) http.Handler

type MiddlewareChain []Middleware

type RouteGroup struct { // Added for versioning
	prefix     string
	router     *Router
	middleware MiddlewareChain
}
var customTypes = map[string]string{}

// NewRouter creates a new Router instance
func NewRouter() *Router {
	return &Router{
		routeGroups: make(map[string]*RouteGroup), // Initialize the routeGroups map
	}
}
func (g *RouteGroup) Handle(method, path string, handler http.HandlerFunc, middleware ...Middleware) {
	fullPath := g.prefix + path // Prepend the group prefix to the path
	g.router.Handle(method, fullPath, handler, middleware...)
}
// Handle adds a new route to the router
func (r *Router) Handle(method, path string, handler http.HandlerFunc, middleware ...Middleware) {
	pattern, params := parsePath(path)
	route := Route{
		path:       path,
		method:     method,
		handler:    handler,
		pattern:    pattern,
		params:     params,
		middleware: middleware, // This now directly assigns the slice of middleware
	}
	r.routes = append(r.routes, route)
}
func (r *Router) Group(version string) *RouteGroup {
	if group, exists := r.routeGroups[version]; exists {
		return group
	}

	group := &RouteGroup{
		prefix: "/api/" + version,
		router: r,
	}
	r.routeGroups[version] = group
	return group
}
// Use adds new middleware to the router
func (r *Router) Use(middleware ...Middleware) {
	r.middleware = append(r.middleware, middleware...) // Appends the new middleware to the existing slice
}

// ServeHTTP makes the router implement the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var finalHandler http.Handler

	finalHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for _, route := range r.routes {
			if matches := route.pattern.FindStringSubmatch(req.URL.Path); matches != nil && route.method == req.Method {
				params := make(map[string]string)
				for i, name := range route.params {
					params[name] = matches[i+1]
				}

				type key int
				const paramsKey key = 0
				ctx := context.WithValue(req.Context(), paramsKey, params)
				req = req.WithContext(ctx)

				routeHandler := http.Handler(route.handler)

				// Applying route-specific middleware in order
				for _, mw := range route.middleware {
					routeHandler = mw(routeHandler)
				}

				routeHandler.ServeHTTP(w, req)
				return
			}
		}
		http.NotFound(w, req)
	})

	// Applying global middleware in order
	for _, mw := range r.middleware {
		finalHandler = mw(finalHandler)
	}

	finalHandler.ServeHTTP(w, req)
}

func AddCustomType(name, pattern string) error {
    if _, exists := customTypes[name]; exists {
        return fmt.Errorf("a custom type with the name '%s' already exists", name)
    }

    // Validate if the provided pattern is a valid regular expression
    _, err := regexp.Compile(pattern)
    if err != nil {
        return fmt.Errorf("invalid pattern: %v", err)
    }

    customTypes[name] = pattern
    return nil
}

func parsePath(path string) (*regexp.Regexp, map[int]string) {
	paramRe := regexp.MustCompile(`{(\w+)(?::(\w+))?}`) // Updated regex to capture optional type

	matches := paramRe.FindAllStringSubmatch(path, -1)
	params := make(map[int]string)

	for i, match := range matches {
		paramName := match[1]
		paramType := match[2]

		// Default regex is \w+ for word characters
		paramRegex := `(\w+)`

		// Customize regex based on parameter type
		switch paramType {
		case "int":
			paramRegex = `(\d+)`
		case "float":
			paramRegex = `(\d+\.\d+)`
		case "uuid":
			paramRegex = `([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})`
		case "alphanumeric":
			paramRegex = `([a-zA-Z0-9]+)`
		case "string":
			paramRegex = `(.+)`
		case "slug":
			paramRegex = `([a-z0-9-]+)`
		case "email":
			paramRegex = `([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`
		case "date":
			paramRegex = `(\d{4}-\d{2}-\d{2})`
		case "ipv4":
			paramRegex = `((?:\d{1,3}\.){3}\d{1,3})`
		case "ipv6":
			paramRegex = `((?:[a-fA-F0-9]{1,4}:){7}[a-fA-F0-9]{1,4})`
		default:
			if customPattern, exists := customTypes[paramType]; exists {
				paramRegex = fmt.Sprintf("(%s)", customPattern)
			}
		}

		path = strings.Replace(path, match[0], paramRegex, 1)
		params[i] = paramName
	}

	pathRe := regexp.MustCompile("^" + path + "$")
	return pathRe, params
}

func UploadFile(r *http.Request, formKey, uploadDir string) (string, error) {
	file, header, err := r.FormFile(formKey)
	if err != nil {
		return "", err
	}
	defer file.Close()

	safeFilename := filepath.Base(header.Filename)
	uploadPath := filepath.Join(uploadDir, safeFilename)
	out, err := os.Create(uploadPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return uploadPath, nil
}
