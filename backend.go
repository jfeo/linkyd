package main

import (
	"errors"
	"net/http"
)

var UnsupportedContentType = errors.New("unsupported content type")

type LinkyBackend interface {
	Create()
	List()
	Delete(id string)
}

func GetBackend(l *Linky, w http.ResponseWriter, r *http.Request) LinkyBackend {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		return NewJSONBackend(l, w, r)
	default:
		return NewHTMXBackend(l, w, r)
	}
}
