package linky

import (
	"log/slog"
	"net/http"
	"text/template"

	"feodor.dk/linkyd/linky/link"
)

type LinkService struct {
	Links link.Repository
}

type LinksAsUser struct {
	AsUser string      `json:"asUser"`
	Links  []link.Link `json:"links"`
}

type Links struct {
	Links []link.Link `json:"links"`
}

func New(links link.Repository) LinkService {
	return LinkService{Links: links}
}

func (l *LinkService) CreateLink(url string, title string, user string) (link.Link, error) {
	if title == "" {
		gottenTitle, err := GetTitleOfLink(url)
		if err == nil {
			title = gottenTitle
		} else {
			title = url
		}
	}

	if user == "" {
		user = "all"
	}

	if ln, err := l.Links.Create(link.Link{
		URL:   url,
		User:  user,
		Title: title,
	}); err != nil {
		slog.Error(
			"Error creating link",
			slog.String("URL", ln.URL),
			slog.String("Title", ln.Title),
			slog.String("User", ln.User),
		)
		return link.Link{}, err
	} else {
		slog.Info("Created link", getLinkSlogArgs(ln)...)
		return ln, nil
	}
}

func (l *LinkService) DeleteLink(id string) (link.Link, error) {
	if ln, err := l.Links.Delete(id); err != nil {
		slog.Info("Error deleting link", slog.String("error", err.Error()), slog.String("ID", id))
		return link.Link{}, err
	} else {
		slog.Info("Deleted link", getLinkSlogArgs(ln)...)
		return ln, nil
	}
}

func (l LinkService) RenderTemplateOr500(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	err := tmpl.Execute(w, l)

	if err != nil {
		r.Response.StatusCode = 500
	}
}

func (l *LinkService) AsUser(user string) LinksAsUser {
	return LinksAsUser{
		AsUser: user,
		Links:  l.Links.AsUser(user),
	}
}

func (l *LinkService) All() Links {
	return Links{Links: l.Links.AsUser("")}
}

func getLinkSlogArgs(ln link.Link) []any {
	args := make([]any, 0)

	args = append(args,
		slog.String("ID", ln.ID),
		slog.String("URL", ln.URL),
		slog.String("Title", ln.Title),
		slog.String("User", ln.User),
		slog.String("Added", ln.AddedAt.String()),
	)
	return args
}
