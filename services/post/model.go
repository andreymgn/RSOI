package post

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

// Post describes a post
type Post struct {
	Uid        string
	Title      string
	URL        string
	CreatedAt  time.Time
	ModifiedAt pq.NullTime
}

type datastore interface {
	getAll(int32, int32) ([]*Post, error)
	getOne(string) (*Post, error)
	create(string, string) error
	update(string, string, string) error
	delete_(string) error
}

type db struct {
	*sql.DB
}

func newDB(connString string) (*db, error) {
	postgres, err := sql.Open("postgres", connString)
	return &db{postgres}, err
}

func (db *db) getAll(pageSize, pageNumber int32) ([]*Post, error) {
	query := "SELECT * FROM posts ORDER BY created_at DESC LIMIT $1, $2"
	lastRecord := pageNumber * pageSize
	rows, err := db.Query(query, lastRecord, pageSize)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	result := make([]*Post, 0)
	for rows.Next() {
		post := new(Post)
		err := rows.Scan(&post.Uid, &post.Title, &post.URL, &post.CreatedAt, &post.ModifiedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *db) getOne(uid string) (*Post, error) {
	query := "SELECT * FROM posts WHERE uid=$1"
	row := db.QueryRow(query, uid)
	result := new(Post)
	if err := row.Scan(&result.Uid, &result.Title, &result.URL, &result.CreatedAt, &result.ModifiedAt); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *db) create(title, url string) error {
	query := "INSERT INTO posts (title, url) VALUES ($1, $2)"
	_, err := db.Exec(query, title, url)
	return err
}

func (db *db) update(title, url, uid string) error {
	query := "UPDATE posts SET title=COALESCE(NULLIF($1,''), title), url=COALESCE(NULLIF($2,''), url) WHERE uid=$3"
	_, err := db.Exec(query, title, url, uid)
	return err
}

func (db *db) delete_(uid string) error {
	query := "DELETE FROM posts WHERE uid=$1"
	_, err := db.Exec(query, uid)
	return err
}
