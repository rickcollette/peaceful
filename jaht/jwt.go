package jaht

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

// GenerateToken generates a new JWT token
func GenerateToken(userID string, expirationTime time.Duration, secretKey []byte) (string, error) {
    // Creating the JWT claims, which includes the user ID and expiry time
    claims := &jwt.StandardClaims{
        Subject:   userID,
        ExpiresAt: time.Now().Add(expirationTime).Unix(),
    }

    // Declaring the token with the algorithm used for signing, and the claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Creating the JWT token
    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// ValidateToken validates the JWT token
func ValidateToken(tokenString string, secretKey []byte) (string, error) {
    claims := &jwt.StandardClaims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })

    if err != nil {
        return "", err
    }

    if !token.Valid {
        return "", errors.New("invalid token")
    }

    return claims.Subject, nil
}

// JwtMiddleware is a middleware function for validating JWT tokens
func JwtMiddleware(next http.Handler, secretKey []byte) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get the JWT token from the Authorization header
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Authorization header must be provided", http.StatusUnauthorized)
            return
        }

        // Validate the JWT token
        _, err := ValidateToken(tokenString, secretKey)
        if err != nil {
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }

        // If the token is valid, call the next handler
        next.ServeHTTP(w, r)
    })
}
