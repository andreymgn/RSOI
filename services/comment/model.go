package comment

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Comment describes comment to a post
type Comment struct {
	UID        uuid.UUID
	PostUID    uuid.UUID
	Body       string
	ParentUID  uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type datastore interface {
	getAll(uuid.UUID, int32, int32) ([]*Comment, error)
	create(uuid.UUID, string, uuid.UUID) (*Comment, error)
	update(uuid.UUID, string) error
	delete(uuid.UUID) error
}

type db struct {
	*sql.DB
}

func newDB(connString string) (*db, error) {
	postgres, err := sql.Open("postgres", connString)
	return &db{postgres}, err
}

func (db *db) getAll(postUID uuid.UUID, pageSize, pageNumber int32) ([]*Comment, error) {
	query := "SELECT * FROM comments ORDER BY created_at DESC LIMIT $1 OFFSET $2"
	lastRecord := pageNumber * pageSize
	rows, err := db.Query(query, pageSize, lastRecord)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	result := make([]*Comment, 0)
	for rows.Next() {
		comment := new(Comment)
		var uid, pUID, parentUID string
		err := rows.Scan(&uid, &pUID, &comment.Body, &parentUID, &comment.CreatedAt, &comment.ModifiedAt)
		if err != nil {
			return nil, err
		}

		comment.UID, err = uuid.Parse(uid)
		if err != nil {
			return nil, err
		}

		comment.PostUID, err = uuid.Parse(pUID)
		if err != nil {
			return nil, err
		}

		if parentUID != "" {
			comment.ParentUID, err = uuid.Parse(parentUID)
			if err != nil {
				return nil, err
			}
		} else {
			comment.ParentUID = uuid.Nil
		}

		result = append(result, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *db) create(postUID uuid.UUID, body string, parentUID uuid.UUID) (*Comment, error) {
	comment := new(Comment)

	query := "INSERT INTO comments (uid, post_uid, body, parent_uid, created_at, modified_at) VALUES ($1, $2, $3, $4, $5, $6)"

	uid := uuid.New()
	now := time.Now()

	comment.UID = uid
	comment.PostUID = postUID
	comment.Body = body
	comment.ParentUID = parentUID
	comment.CreatedAt = now
	comment.ModifiedAt = now

	_, err := db.Query(query, uid.String(), postUID.String(), body, parentUID.String(), now, now)
	return comment, err
}

func (db *db) update(uid uuid.UUID, body string) error {
	query := "UPDATE comments SET body=$1, modified_at=$2 WHERE uid=$3"
	_, err := db.Exec(query, body, time.Now(), uid.String())
	return err
}

func (db *db) delete(uid uuid.UUID) error {
	query := "DELETE FROM POSTS WHERE uid=$1"
	_, err := db.Exec(query, uid.String())
	return err
}
