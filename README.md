# peaceful
This is a lightweight (as of now) RESTful library for go.

This library provides a collection of utilities for building web applications with Go. It includes a router and various middleware for tasks like caching, CSRF protection, CORS handling, content negotiation, and request binding.

## Why?
  - Current implementations are more complex than needed
  - Some of the super popular ones have been sidelined or are no longer maintained
  - I wanted a lightweight (as of now) RESTful library for go

## Support

Please file an issue in github for this project.  

This includes:
  - Feature requests
  - Contributions
  - Bugs

## Installation

go get github.com/rickcollette/peaceful/router

## Usage

### Router

peaceful router provides a basic router for handling HTTP requests. Here's how you can use it:

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/rickcollette/peaceful/router"  // Replace with the actual path
)

func main() {
    r := router.NewRouter()

    r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, world!")
    })

    http.ListenAndServe(":8080", r)
}
```

### Caching

Use peaceful router's caching middleware to cache HTTP GET requests. The corrected way to add the caching middleware is shown below.

```go
r := router.NewRouter()

r.Use(func(next http.Handler) http.Handler {
    return router.CachingMiddleware(10 * time.Minute, next)
})

// Your routes here

```

### Shortcuts

peaceful router provides shortcut methods for common HTTP methods like GET, POST, PUT, DELETE. They are used like this:

```go
r.GET("/path", handlerFunc)
r.POST("/path", handlerFunc)
r.PUT("/path", handlerFunc)
r.DELETE("/path", handlerFunc)
```

### Request Binding

peaceful router contains functions for binding request data to structs, including JSON and XML data. Example:

```go
type MyData struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

var data MyData
err := router.BindJSON(r, &data)
if err != nil {
    // Handle error
}
```

### CSRF Protection

peaceful router provides CSRF protection middleware. Use it like this:

```go
r.Use(router.CSRFMiddleware)
```

### CORS Handling

peaceful router provides CORS handling middleware with configurable options. Here’s an example of how to use it:

```go
options := router.CORSOptions{
    AllowedOrigins: []string{"*"},  // Allow all origins
    AllowedMethods: []string{"GET", "POST"},  // Allow only GET and POST requests
}

r.Use(router.CORS(options))
```

### Content Negotiation

peaceful router provides a function for handling content negotiation. Here's how to use it:

```go
router.Respond(w, r, 200, data)  // Automatically selects the content type based on the "Accept" header
```

## Example RESTful Application

Here is a complete example of a RESTful application that utilizes peaceful:

```go
package main

import (
    "net/http"
    "time"
    "github.com/rickcollette/peaceful/router"  
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"message": "Hello, world!"}
    router.Respond(w, r, 200, data)
}

func main() {
    r := router.NewRouter()

    // Corrected Middleware - Wrapping CachingMiddleware in a function that matches router.Middleware type
    r.Use(func(next http.Handler) http.Handler {
        return router.CachingMiddleware(10*time.Minute, next)
    })

    r.Use(router.CSRFMiddleware)

    options := router.CORSOptions{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST"},
    }
    r.Use(router.CORS(options))

    // Routes
    r.GET("/hello", helloHandler)

    // Start server
    http.ListenAndServe(":8080", r)
}
```
