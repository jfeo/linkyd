package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"golang.org/x/net/html"
)

type Link struct {
	ID      string    `json:"id"`
	URL     string    `json:"url"`
	Title   string    `json:"title"`
	AddedAt time.Time `json:"addedAt"`
}

type Linky struct {
	NextID int             `json:"nextID"`
	Links  map[string]Link `json:"links"`
}

func GetTitleOfLink(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("could not access link, got status %d", resp.StatusCode)
	}

	tree, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not parse link: %w", err)
	}

	if title, err := getTitleFromHTML(tree); err != nil {
		return "", err
	} else {
		return title, nil
	}

}

var TitleNotFound = fmt.Errorf("title was not found in html tree")

func getTitleFromHTML(tree *html.Node) (string, error) {
	if tree == nil {
		return "", TitleNotFound
	} else if tree.Type == html.ElementNode && tree.Data == "title" {
		return getAllSiblingText(tree.FirstChild), nil
	} else if title, err := getTitleFromHTML(tree.FirstChild); err == nil {
		return title, nil
	} else if title, err := getTitleFromHTML(tree.NextSibling); err == nil {
		return title, nil
	} else {
		return "", err
	}
}

func getAllSiblingText(sibling *html.Node) string {
	if sibling == nil {
		return ""
	} else if sibling.Type != html.TextNode {
		return getAllSiblingText(sibling.FirstChild) + getAllSiblingText(sibling.NextSibling)
	} else {
		return sibling.Data
	}
}

func (l *Linky) CreateLink(url string, title string) Link {
	id := strconv.Itoa(l.NextID)
	if title == "" {
		gottenTitle, err := GetTitleOfLink(url)
		if err == nil {
			title = gottenTitle
		}
	}

	l.Links[id] = Link{ID: id, URL: url, Title: title, AddedAt: time.Now()}
	l.NextID++
	slog.Info("Added link", "ID", id, "URL", url, "Title", title)

	return l.Links[id]
}

func (l *Linky) DeleteLink(id string) Link {
	slog.Info("Removing link", "ID", id)
	link := l.Links[id]
	delete(l.Links, id)
	return link
}

func (l Linky) RenderTemplateOr500(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	err := tmpl.Execute(w, l)

	if err != nil {
		r.Response.StatusCode = 500
	}
}

func (l Linky) SaveLinks() {
	fd, err := os.Create("linky.json")
	if err != nil {
		slog.Error("Could not save links", "error", err)
		return
	}
	defer fd.Close()

	marshalledLinks, err := json.Marshal(l)
	if err != nil {
		slog.Error("Could not save links", "error", err)
		return
	}
	fd.Write(marshalledLinks)
}
