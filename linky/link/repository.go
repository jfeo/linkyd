package link

import (
	"errors"
)

type Repository interface {
	AsUser(user string) []Link
	Create(link Link) (Link, error)
	Delete(linkID string) (Link, error)
}

var ErrLinkNotFound error = errors.New("link not found")
