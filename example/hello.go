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
