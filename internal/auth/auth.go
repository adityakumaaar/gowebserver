package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Get the API Key from the Headers
// Authorization: ApiKey <value>
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("No auth string present in the headers")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed ApiKey")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed ApiKey")
	}
	return vals[1], nil
}
