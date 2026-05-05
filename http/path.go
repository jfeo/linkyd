package http

import (
	"net/http"
	"strings"
)

func getPathSegment(r *http.Request, argumentIndex int) (string, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != argumentIndex+1 {
		return "", ErrInvalidPath
	}

	return parts[argumentIndex], nil
}
