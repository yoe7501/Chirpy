package auth

import (
	"errors"
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	h := headers.Values("Authorization")
	if len(h) == 0 {
		return "", errors.New("no authorization header")
	}
	if len(h[0]) <= 7 {
		return "", errors.New("malformed Authorization header")
	}
	return h[0][7:], nil
}
