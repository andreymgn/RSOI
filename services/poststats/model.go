package poststats

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	errNotFound   = errors.New("post statustics not found")
	errNotCreated = errors.New("post stats not created")
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
	switch err := row.Scan(&uidString, &result.NumLikes, &result.NumDislikes, &result.NumViews); err {
	case nil:
		result.UID = uid
		return result, nil
	case sql.ErrNoRows:
		return nil, errNotFound
	default:
		return nil, err
	}
}

func (db *db) create(uid uuid.UUID) (*PostStats, error) {
	query := "INSERT INTO posts_stats (post_uid, num_likes, num_dislikes, num_views) VALUES ($1, 0, 0, 0)"
	ps := new(PostStats)
	ps.UID = uid
	ps.NumLikes = 0
	ps.NumDislikes = 0
	ps.NumViews = 0
	result, err := db.Exec(query, uid.String())
	nRows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if nRows == 0 {
		return nil, errNotCreated
	}

	return ps, nil
}

func (db *db) like(uid uuid.UUID) error {
	query := "UPDATE posts_stats SET num_likes = num_likes + 1 WHERE post_uid=$1"
	result, err := db.Exec(query, uid.String())
	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if nRows == 0 {
		return errNotFound
	}

	return nil
}

func (db *db) dislike(uid uuid.UUID) error {
	query := "UPDATE posts_stats SET num_dislikes = num_dislikes + 1 WHERE post_uid=$1"
	result, err := db.Exec(query, uid.String())
	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if nRows == 0 {
		return errNotFound
	}

	return nil
}

func (db *db) view(uid uuid.UUID) error {
	query := "UPDATE posts_stats SET num_views = num_views + 1 WHERE post_uid=$1"
	result, err := db.Exec(query, uid.String())
	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if nRows == 0 {
		return errNotFound
	}

	return nil
}

func (db *db) delete(uid uuid.UUID) error {
	query := "DELETE FROM posts_stats WHERE post_uid=$1"
	result, err := db.Exec(query, uid.String())
	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if nRows == 0 {
		return errNotFound
	}

	return nil
}
