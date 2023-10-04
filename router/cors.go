package router

import (
	"net/http"
	"strings"
)

type CORSOptions struct {
    AllowedOrigins   []string
    AllowedMethods   []string
    AllowedHeaders   []string
    AllowCredentials bool
    ExposeHeaders    []string
}

func CORS(opts CORSOptions) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Set headers based on CORS options
            if len(opts.AllowedOrigins) > 0 {
                w.Header().Set("Access-Control-Allow-Origin", strings.Join(opts.AllowedOrigins, ", "))
            }

            if len(opts.AllowedMethods) > 0 {
                w.Header().Set("Access-Control-Allow-Methods", strings.Join(opts.AllowedMethods, ", "))
            }

            if len(opts.AllowedHeaders) > 0 {
                w.Header().Set("Access-Control-Allow-Headers", strings.Join(opts.AllowedHeaders, ", "))
            }

            if opts.AllowCredentials {
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }

            if len(opts.ExposeHeaders) > 0 {
                w.Header().Set("Access-Control-Expose-Headers", strings.Join(opts.ExposeHeaders, ", "))
            }

            // Handle preflight requests
            if r.Method == "OPTIONS" {
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}