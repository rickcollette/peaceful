package router

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

type contextKey string

const (
	requestIDKey   contextKey = "requestID"
	contentTypeKey contextKey = "content-type"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mtx      sync.Mutex
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL, time.Since(start))
	})
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getVisitor(ip string) *rate.Limiter {
	mtx.Lock()
	defer mtx.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(1, 5) // adjust as needed
		visitors[ip] = limiter
	}

	return limiter
}

func getIP(r *http.Request) string {
	// Extract the X-Forwarded-For header value
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// If the header is present, extract the first IP address from the list
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	// Otherwise, return the direct client IP
	return r.RemoteAddr
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getIP(r)
		limiter := getVisitor(clientIP)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ContentNegotiationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptHeader := r.Header.Get("Accept")
		var contentType string

		switch {
		case strings.Contains(acceptHeader, "application/json"):
			contentType = "application/json"
		case strings.Contains(acceptHeader, "application/xml"):
			contentType = "application/xml"
		default:
			contentType = "text/html" // or any default content type
		}

		ctx := context.WithValue(r.Context(), contentTypeKey, contentType)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
