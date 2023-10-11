package router

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Router struct {
	routes      []Route
	middleware  MiddlewareChain // Changed to MiddlewareChain
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

var (
	customTypes        = map[string]string{}
	validFilenameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+\.[a-zA-Z0-9]+$`)
)

var typePatterns = map[string]*regexp.Regexp{
	"int":          regexp.MustCompile(`\d+`),
	"float":        regexp.MustCompile(`\d+(\.\d+)?`),
	"uuid":         regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`),
	"alphanumeric": regexp.MustCompile(`[a-zA-Z0-9]+`),
	"string":       regexp.MustCompile(`[^/]+`),
	"slug":         regexp.MustCompile(`[a-z0-9]+(-[a-z0-9]+)*`),
	"email":        regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
	"date":         regexp.MustCompile(`\d{4}-\d{2}-\d{2}`),
	"ipv4":         regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`),
	"ipv6":         regexp.MustCompile(`\b([a-fA-F0-9]{1,4}:){7}[a-fA-F0-9]{1,4}\b`),
}

func isValidFilename(filename string) bool {
	return validFilenameRegex.MatchString(filename)
}

// NewRouter creates a new Router instance
func NewRouter() *Router {
	return &Router{
		routeGroups: make(map[string]*RouteGroup), // Initialize the routeGroups map
	}
}
func (g *RouteGroup) Handle(method, path string, handler http.HandlerFunc, middleware ...Middleware) {
	fullPath := g.prefix + path // Prepend the group prefix to the path
	middleware = append(middleware, g.middleware...) // Add the group's middleware to the route's middleware
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
// Modified ServeHTTP function to apply middleware
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

				// If you want to apply global middleware to each route as well,
				// you can do so here by iterating over r.middleware.

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

// Param function to extract parameters from the request context
func Param(r *http.Request, name string) string {
	type key int
	const paramsKey key = 0
	params := r.Context().Value(paramsKey).(map[string]string)
	return params[name]
}

// RegisterParamType allows the registration of custom parameter types with specific regex patterns
func RegisterParamType(name, pattern string) {
	typePatterns[name] = regexp.MustCompile(pattern)
}

func parsePath(path string) (*regexp.Regexp, map[int]string) {
	paramRe := regexp.MustCompile(`{(\w+)(?::(\w+))?}`) // Updated regex to capture optional type

	matches := paramRe.FindAllStringSubmatch(path, -1)
	params := make(map[int]string)

	for i, match := range matches {
		paramName := match[1]
		paramType := match[2]

		// Default regex is \w+ for word characters
		paramRegex := typePatterns["string"].String() // Default to string type

		// If the parameter type is recognized, use its specific regex pattern
		if pattern, exists := typePatterns[paramType]; exists {
			paramRegex = pattern.String()
		}

		path = strings.Replace(path, match[0], fmt.Sprintf("(%s)", paramRegex), 1)
		params[i] = paramName
	}

	pathRe := regexp.MustCompile("^" + path + "$")
	return pathRe, params
}


func UploadFile(r *http.Request, formKey, uploadDir string) (string, error) {
	// Parse the form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		return "", err
	}

	// Get the file from the form data
	file, header, err := r.FormFile(formKey)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Validate the filename
	if !isValidFilename(header.Filename) {
		return "", errors.New("invalid filename")
	}

	// Construct the full path safely
	safeFilename := filepath.Base(header.Filename)
	uploadPath := filepath.Join(uploadDir, safeFilename)

	// Create the output file
	out, err := os.Create(uploadPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy the uploaded file to the output file
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return uploadPath, nil
}