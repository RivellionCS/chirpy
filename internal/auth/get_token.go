package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	strippedHeader := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	if strippedHeader == "" {
		return "", fmt.Errorf("Error empty token")
	}
	return strippedHeader, nil
}