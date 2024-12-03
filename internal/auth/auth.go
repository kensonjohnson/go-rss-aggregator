package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GETAPIKey extracts an API Key from the headers of a HTTP request
//
// Example:
//
// "Authorization: ApiKey {user's apikey here}"
func GetAPIKey(headers http.Header) (string, error) {
	value := headers.Get("Authorization")
	if value == "" {
		return "", errors.New("authorization header not set")
	}

	parts := strings.Split(value, " ")
	if len(parts) != 2 {
		return "", errors.New("malformed Authorization header")
	}

	if parts[0] != "ApiKey" {
		return "", errors.New("incorrect key format specified")
	}

	if len(parts[1]) != 64 {
		return "", errors.New("ApiKey must be exactly 64 bytes")
	}

	return parts[1], nil
}
