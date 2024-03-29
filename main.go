package main

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

func main() {
	linky := Linky{Links: make(map[string]Link)}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b := GetBackend(&linky, w, r)
		b.List()
	})

	http.HandleFunc("/as/", func(w http.ResponseWriter, r *http.Request) {
		b := GetBackend(&linky, w, r)
		if asUser, err := getPathSegment(r, 2); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			b.As(asUser)
		}
	})

	http.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
		b := GetBackend(&linky, w, r)

		switch r.Method {
		case "POST":
			b.Create()
		case "GET":
			b.List()
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/links/", func(w http.ResponseWriter, r *http.Request) {
		b := GetBackend(&linky, w, r)

		if r.Method != "DELETE" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		id, err := getPathSegment(r, 2)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b.Delete(id)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getPathSegment(r *http.Request, argumentIndex int) (string, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != argumentIndex+1 {
		return "", errors.New("Invalid path")
	}

	return parts[argumentIndex], nil
}
