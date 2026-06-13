package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("Header is empty")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("Authorization header must start with 'Bearer '")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	if token == "" {
		return "", fmt.Errorf("Error empty token")
	}

	return token, nil
}