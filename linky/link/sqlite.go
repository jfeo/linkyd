package link

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type SQLiteLinkRepository struct {
	db *sql.DB
}

const (
	createTable = `
		CREATE TABLE IF NOT EXISTS links (
			id INTEGER,
			uuid BINARY(16),
			title VARCHAR(256),
			url VARCHAR(512),
			user VARCHAR(32),
			addedAt DATETIME,
			PRIMARY KEY (id),
			UNIQUE(uuid)
		);
	`
	asUserQuery = `SELECT uuid, title, url, user, addedAt FROM links WHERE user != ? ORDER BY addedAt;`
	insertQuery = `INSERT INTO links (uuid, title, url, user, addedAt) VALUES (?, ?, ?, ?, ?);`
	deleteQuery = `DELETE FROM links WHERE uuid = ?;`
)

func NewSQLiteLinkRepository() (*SQLiteLinkRepository, error) {
	if db, err := sql.Open("sqlite3", "db.sqlite"); err != nil {
		slog.Error("error opening sqlite3 db", slog.String("error", err.Error()))
		return nil, err
	} else if _, err := db.Exec(createTable); err != nil {
		slog.Error("error creating schema", slog.String("error", err.Error()))
		return nil, err
	} else {
		return &SQLiteLinkRepository{db: db}, nil
	}
}

func (r *SQLiteLinkRepository) AsUser(user string) []Link {
	links := make([]Link, 0)

	if rows, err := r.db.Query(asUserQuery, user); err != nil {
		slog.Error("error querying links as user", slog.String("error", err.Error()), slog.String("user", user))
	} else {
		for rows.Next() {
			if link, err := scanRow(rows); err != nil {
				slog.Error("error scanning row", slog.String("error", err.Error()))
			} else {
				links = append(links, link)
			}
		}
	}

	return links
}

func (r *SQLiteLinkRepository) Create(link Link) (Link, error) {
	if link.AddedAt.IsZero() {
		link.AddedAt = time.Now()
	}

	if linkUuid, err := uuid.NewRandom(); err != nil {
		return link, err
	} else if _, err := r.db.Exec(insertQuery, linkUuid, link.Title, link.URL, link.User, link.AddedAt); err != nil {
		return link, err
	} else {
		link.ID = linkUuid.String()
		return link, nil
	}
}

func (r *SQLiteLinkRepository) Delete(linkID string) (Link, error) {
	if linkUuid, err := uuid.Parse(linkID); err != nil {
		return Link{}, err
	} else if _, err := r.db.Exec(deleteQuery, linkUuid); err != nil {
		return Link{}, err
	} else {
		return Link{}, nil
	}
}

func scanRow(rows *sql.Rows) (Link, error) {
	var link Link
	err := rows.Scan(&link.ID, &link.Title, &link.URL, &link.User, &link.AddedAt)
	return link, err
}
