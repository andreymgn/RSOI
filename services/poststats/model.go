package poststats

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// PostStats describes post statistics
type PostStats struct {
	Uid         string
	NumLikes    int32
	NumDislikes int32
	NumViews    int32
}

type datastore interface {
	get(string) (*PostStats, error)
	create(string) error
	like(string) error
	dislike(string) error
	view(string) error
	delete(string) error
}

type db struct {
	*sql.DB
}

func newDB(connString string) (*db, error) {
	postgres, err := sql.Open("postgres", connString)
	return &db{postgres}, err
}

func (db *db) get(uid string) (*PostStats, error) {
	query := "SELECT * FROM posts_stats WHERE post_uid=$1"
	row := db.QueryRow(query, uid)
	result := new(PostStats)
	if err := row.Scan(&result.Uid, &result.NumLikes, &result.NumDislikes, &result.NumViews); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *db) create(uid string) error {
	query := "INSERT INTO posts_stats (post_uid) VALUES ($1)"
	_, err := db.Query(query, uid)
	return err
}

func (db *db) like(uid string) error {
	query := "UPDATE posts_stats SET num_likes = num_likes + 1 WHERE post_uid=$1"
	_, err := db.Exec(query, uid)
	return err
}

func (db *db) dislike(uid string) error {
	query := "UPDATE posts_stats SET num_dislikes = num_dislikes + 1 WHERE post_uid=$1"
	_, err := db.Exec(query, uid)
	return err
}

func (db *db) view(uid string) error {
	query := "UPDATE posts_stats SET num_views = num_views + 1 WHERE post_uid=$1"
	_, err := db.Exec(query, uid)
	return err
}

func (db *db) delete(uid string) error {
	query := "DELETE FROM posts_stats WHERE post_uid=$1"
	_, err := db.Exec(query, uid)
	return err
}
