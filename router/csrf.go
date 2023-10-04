package router

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/go-playground/validator/v10" // Importing a third-party validation library
)

// CSRF Middleware
func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("csrf_token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		csrfToken := r.Header.Get("X-CSRF-Token")
		if csrfToken == "" {
			csrfToken = r.PostFormValue("csrf_token")
		}

		if csrfToken != cookie.Value {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Function to generate and set CSRF token cookie
func SetCSRFToken(w http.ResponseWriter) {
	token := make([]byte, 32)
	rand.Read(token)
	encodedToken := base64.StdEncoding.EncodeToString(token)

	http.SetCookie(w, &http.Cookie{
		Name:  "csrf_token",
		Value: encodedToken,
		Path:  "/",
	})
}

// Validation function using a third-party library
var validate = validator.New()

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
