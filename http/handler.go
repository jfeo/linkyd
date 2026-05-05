package http

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"feodor.dk/linkyd/backend"
	"feodor.dk/linkyd/linky"
	"feodor.dk/linkyd/static"
)

var ErrInvalidPath = errors.New("invalid path")

func ListenAndServe(linky *linky.LinkService, port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b := backend.Get(linky, w, r)
		b.List()
	})

	http.HandleFunc("/as/", func(w http.ResponseWriter, r *http.Request) {
		b := backend.Get(linky, w, r)
		if asUser, err := getPathSegment(r, 2); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			b.As(asUser)
		}
	})

	http.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
		b := backend.Get(linky, w, r)

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
		b := backend.Get(linky, w, r)

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

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Add("Content-Type", "image/x-icon")
		w.Header().Add("Content-Length", fmt.Sprintf("%d", len(static.Favicon)))
		w.Write(static.Favicon)
	})

	listenAddr := fmt.Sprintf(":%d", port)
	slog.Info("Starting HTTP listener", slog.String("address", listenAddr))
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
