package linky

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

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

var ErrTitleNotFound = errors.New("title was not found in html tree")

func getTitleFromHTML(tree *html.Node) (string, error) {
	if tree == nil {
		return "", ErrTitleNotFound
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
