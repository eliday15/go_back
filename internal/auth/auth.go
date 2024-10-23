package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Extracts the API key from the headers of an HTTP request
// Example: Authorization: ApiKey {api_key}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authorization header found")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed authorization header")
	}
	apiKey := vals[1]
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}
	return apiKey, nil
}
