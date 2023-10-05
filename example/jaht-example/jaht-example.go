package main

import (
	"net/http"
	"time"
	"github.com/rickcollette/peaceful/jaht"
	"github.com/rickcollette/peaceful/router"
)

func main() {
	// Create a new router
	r := router.NewRouter()

	// Secret key for signing JWT tokens
	secretKey := []byte("your-secret-key")

	// Middleware to validate JWT tokens
	// Corrected the middleware usage here
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
