package link

import (
	"fmt"
	"slices"
	"sync"
	"time"
)

type InMemoryLinkRepository struct {
	lock   sync.Mutex
	nextID uint64
	links  map[string]Link
}

func NewInMemoryLinkRepository() *InMemoryLinkRepository {
	return &InMemoryLinkRepository{
		links: make(map[string]Link),
	}
}

func (r *InMemoryLinkRepository) incrementID() uint64 {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.nextID += 1
	return r.nextID
}

func (r *InMemoryLinkRepository) AsUser(user string) []Link {
	links := make([]Link, 0)

	for _, l := range r.links {
		if l.User != user {
			links = append(links, l)
		}
	}

	slices.SortStableFunc(links, func(a, b Link) int {
		return a.AddedAt.Compare(b.AddedAt)
	})

	return links
}

func (r *InMemoryLinkRepository) Create(link Link) (Link, error) {
	id := r.incrementID()
	r.lock.Lock()
	defer r.lock.Unlock()

	link.ID = fmt.Sprintf("%d", id)
	if link.AddedAt.IsZero() {
		link.AddedAt = time.Now()
	}

	r.links[link.ID] = link

	return link, nil
}

func (r *InMemoryLinkRepository) Delete(linkID string) (Link, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if link, found := r.links[linkID]; found {
		delete(r.links, linkID)
		return link, nil
	} else {
		return Link{}, ErrLinkNotFound
	}
}
