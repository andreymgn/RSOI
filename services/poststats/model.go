package poststats

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// PostStats describes post statistics
type PostStats struct {
	UID         uuid.UUID
	NumLikes    int32
	NumDislikes int32
	NumViews    int32
}

type datastore interface {
	get(uuid.UUID) (*PostStats, error)
	create(uuid.UUID) (*PostStats, error)
	like(uuid.UUID) error
	dislike(uuid.UUID) error
	view(uuid.UUID) error
	delete(uuid.UUID) error
}

type db struct {
	*sql.DB
}

func newDB(connString string) (*db, error) {
	postgres, err := sql.Open("postgres", connString)
	return &db{postgres}, err
}

func (db *db) get(uid uuid.UUID) (*PostStats, error) {
	query := "SELECT * FROM posts_stats WHERE post_uid=$1"
	row := db.QueryRow(query, uid.String())
	result := new(PostStats)
	var uidString string
	if err := row.Scan(&uidString, &result.NumLikes, &result.NumDislikes, &result.NumViews); err != nil {
		return nil, err
	}
	result.UID = uid

	return result, nil
}

func (db *db) create(uid uuid.UUID) (*PostStats, error) {
	query := "INSERT INTO posts_stats (post_uid, num_likes, num_dislikes, num_views) VALUES ($1, 0, 0, 0)"
	result := new(PostStats)
	result.UID = uid
	result.NumLikes = 0
	result.NumDislikes = 0
	result.NumViews = 0
	_, err := db.Query(query, uid.String())
	return result, err
}

func (db *db) like(uid uuid.UUID) error {
	query := "UPDATE posts_stats SET num_likes = num_likes + 1 WHERE post_uid=$1"
	_, err := db.Exec(query, uid.String())
	return err
}

func (db *db) dislike(uid uuid.UUID) error {
	query := "UPDATE posts_stats SET num_dislikes = num_dislikes + 1 WHERE post_uid=$1"
	_, err := db.Exec(query, uid.String())
	return err
}

func (db *db) view(uid uuid.UUID) error {
	query := "UPDATE posts_stats SET num_views = num_views + 1 WHERE post_uid=$1"
	_, err := db.Exec(query, uid.String())
	return err
}

func (db *db) delete(uid uuid.UUID) error {
	query := "DELETE FROM posts_stats WHERE post_uid=$1"
	_, err := db.Exec(query, uid.String())
	return err
}
