# peaceful
This is a lightweight (as of now) set of RESTful libraries for go.

This library provides a collection of utilities for building web applications with Go. It includes a router and various middleware for tasks like caching, CSRF protection, CORS handling, content negotiation, jwt handling, and request binding.

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

## Enhanced Parameterized Routes with Type Support
Peaceful router supports typed parameters in routes, which allows for more granular control over the URL parameters and ensures that the incoming parameters adhere to the expected format. This feature is particularly useful for validating and filtering input data at the routing level.

## Supported Parameter Types
Peaceful router supports various parameter types, including but not limited to:

  - int: Matches integer values, e.g., 42.
  - float: Matches floating-point numbers, e.g., 3.14.
  - uuid: Matches UUIDs, e.g., 123e4567-e89b-12d3-a456-426614174000.
  - alphanumeric: Matches alphanumeric strings, e.g., abc123.
  - string: Matches any string, e.g., hello.
  - slug: Matches URL slugs, e.g., hello-world.
  - email: Matches email addresses, e.g., example@example.com.
  - date: Matches date strings in YYYY-MM-DD format, e.g., 2023-10-04.
  - ipv4: Matches IPv4 addresses, e.g., 192.168.0.1.
  - ipv6: Matches IPv6 addresses, e.g., 2001:0db8:85a3:0000:0000:8a2e:0370:7334.

## Using Typed Parameters in Routes
You can specify the expected type of a parameter directly in the route definition. Peaceful router will then use the appropriate regular expression to match and validate the parameter.


```go
package main

import (
    "fmt"
    "net/http"
    "github.com/rickcollette/peaceful/router"
)

func main() {
    r := router.NewRouter()

    // Integer parameter
    r.GET("/users/:id:int", func(w http.ResponseWriter, r *http.Request) {
        id := router.Param(r, "id")
        fmt.Fprintf(w, "User ID (integer): %s", id)
    })

    // UUID parameter
    r.GET("/products/:uuid:uuid", func(w http.ResponseWriter, r *http.Request) {
        uuid := router.Param(r, "uuid")
        fmt.Fprintf(w, "Product UUID: %s", uuid)
    })

    // Date parameter
    r.GET("/events/:date:date", func(w http.ResponseWriter, r *http.Request) {
        date := router.Param(r, "date")
        fmt.Fprintf(w, "Event Date: %s", date)
    })

    http.ListenAndServe(":8080", r)
}
```
Explanation:
In this updated example:

  - The route ```/users/:id:int``` expects an integer parameter named id.
  - The route ```/products/:uuid:uuid``` expects a UUID parameter named uuid.
  - The route ```/events/:date:date``` expects a date parameter named date.
  - The ```router.Param``` function is used to extract the parameter value from the request.

## Custom Parameter Types
Peaceful also allows the definition of custom parameter types. You can add your own regular expressions to match specific patterns tailored to your application's needs.

Example:
```go
// Register a custom parameter type for matching RGB color codes
router.RegisterParamType("color", `[a-fA-F0-9]{6}`)

r.GET("/colors/:color:color", func(w http.ResponseWriter, r *http.Request) {
    color := router.Param(r, "color")
    fmt.Fprintf(w, "Color Code: #%s", color)
})
```

This flexibility ensures that your application can have robust and dynamic routing capable of handling a wide variety of scenarios, making the development process more efficient and the application more secure and user-friendly.

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

## Implementing JWT Authentication with JAHT

JAHT is a package in the peaceful library that aids in implementing JWT (JSON Web Tokens) authentication in your Go applications. It provides utility functions for generating, validating, and parsing JWTs.

### Installation

To install JAHT, you can use the go get command as shown below:

```sh
go get github.com/rickcollette/peaceful/jaht
```

### Usage

Here's a basic example of how you can use JAHT to implement JWT authentication:

```go
package main

import (
    "net/http"
    "time"
    "github.com/rickcollette/peaceful/jaht"
    "github.com/rickcollette/peaceful/router"
)

func main() {
    r := router.NewRouter()

    // Secret key for signing JWT tokens
    secretKey := []byte("your-secret-key")

    // Middleware to validate JWT tokens
    r.Use(func(next http.Handler) http.Handler {
        return jaht.JwtMiddleware(next, secretKey)
    })

    // Route to generate a JWT token
    r.GET("/generate-token", func(w http.ResponseWriter, r *http.Request) {
        userID := "123"  // Replace with actual user ID
        expirationTime := time.Hour * 24  // Token expiration time

        // Generate JWT token
        token, err := jaht.GenerateToken(userID, expirationTime, secretKey)
        if err != nil {
            http.Error(w, "Failed to generate token", http.StatusInternalServerError)
            return
        }

        data := map[string]string{"token": token}
        router.Respond(w, r, http.StatusOK, data)
    })

    // Protected route that requires a valid JWT token
    r.GET("/protected", func(w http.ResponseWriter, r *http.Request) {
        data := map[string]string{"message": "Welcome to the protected route!"}
        router.Respond(w, r, http.StatusOK, data)
    })

    // Start the server
    http.ListenAndServe(":8080", r)
}
```

In this example, the `JwtMiddleware` is added to the router to validate JWT tokens on all incoming requests. A `/generate-token` route is added to generate a JWT token, and a protected `/protected` route requires a valid JWT token to access.

### Testing JWT Authentication

You can test JWT authentication using curl commands.

1. Generate a JWT token:

```sh
curl http://localhost:8080/generate-token
```

2. Use the generated JWT token to access the protected route:

```sh
curl -H "Authorization: your-jwt-token" http://localhost:8080/protected
```

Replace "your-jwt-token" with the actual token received from the `/generate-token` endpoint.
