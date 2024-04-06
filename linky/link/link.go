package link

import "time"

type Link struct {
	ID      string    `json:"id"`
	URL     string    `json:"url"`
	Title   string    `json:"title"`
	User    string    `json:"user"`
	AddedAt time.Time `json:"addedAt"`
}
