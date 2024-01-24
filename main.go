package main

import (
	"log"
	"net/http"
)

func main() {
	linky := Linky{Links: make(map[string]Link)}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b := GetBackend(&linky, w, r)
		b.List()
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
			w.WriteHeader(405)
			return
		}

		id, err := getPathId(r, 2)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		b.Delete(id)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
