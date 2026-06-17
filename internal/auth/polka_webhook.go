package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiHeader := headers.Get("Authorization")

	if apiHeader == "" {
		return "", fmt.Errorf("Header is empty")
	}

	if !strings.HasPrefix(apiHeader, " ApiKey ") {
		return "", fmt.Errorf("Header must start with ' ApiKey '")
	}

	apiKey := strings.TrimSpace(strings.TrimPrefix(apiHeader, " ApiKey "))

	if apiKey == "" {
		return "", fmt.Errorf("ApiKey is empty")
	}

	return apiKey, nil
}