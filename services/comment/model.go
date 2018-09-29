package comment

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// Comment describes comment to a post
type Comment struct {
	Uid        string
	PostUid    string
	Body       string
	ParentUid  string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type datastore interface {
	getAll(string, int32, int32) ([]*Comment, error)
	create(string, string, string) error
	update(string, string) error
	delete(string) error
}

type db struct {
	*sql.DB
}

func newDB(connString string) (*db, error) {
	postgres, err := sql.Open("postgres", connString)
	return &db{postgres}, err
}

func (db *db) getAll(postUid string, pageSize, pageNumber int32) ([]*Comment, error) {
	query := "SELECT * FROM comments ORDER BY created_at DESC LIMIT $1, $2"
	lastRecord := pageNumber * pageSize
	rows, err := db.Query(query, lastRecord, pageSize)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	result := make([]*Comment, 0)
	for rows.Next() {
		comment := new(Comment)
		err := rows.Scan(&comment.Uid, &comment.PostUid, &comment.Body, &comment.ParentUid, &comment.CreatedAt, &comment.ModifiedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *db) create(postUid, body, parentUid string) error {
	query := "INSERT INTO comments (post_uid, body, parent_uid) VALUES ($1, $2, $3)"
	_, err := db.Query(query, postUid, body, parentUid)
	return err
}

func (db *db) update(uid, body string) error {
	query := "UPDATE comments SET body=$1 WHERE uid=$2"
	_, err := db.Exec(query, body, uid)
	return err
}

func (db *db) delete(uid string) error {
	query := "DELETE FROM POSTS WHERE uid=$1"
	_, err := db.Exec(query, uid)
	return err
}
