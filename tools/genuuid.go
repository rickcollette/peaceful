package tools

import "github.com/google/uuid"

// GenerateToken generates a new JWT token
func GenerateUUID() (string, error) {
    uuid := uuid.New()
    uuidString := uuid.String()

    return uuidString, nil
}

