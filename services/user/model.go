package user

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	errNotFound           = errors.New("user not found")
	errNotCreated         = errors.New("user not created")
	errUserPostNotCreated = errors.New("user post not created")
)

const (
	bcryptCost = 14
)

// User describes public user info
type User struct {
	UID      uuid.UUID
	Username string
}

type datastore interface {
	getUserInfo(uuid.UUID) (*User, error)
	create(string, string) (*User, error)
	update(uuid.UUID, string) error
	delete(uuid.UUID) error
	checkPassword(uuid.UUID, string) (bool, error)
	getUIDByUsername(string) (uuid.UUID, error)
}

type db struct {
	*sql.DB
}

func newDB(connString string) (*db, error) {
	postgres, err := sql.Open("postgres", connString)
	return &db{postgres}, err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (db *db) getUserInfo(uid uuid.UUID) (*User, error) {
	query := "SELECT username FROM users WHERE uid=$1"
	row := db.QueryRow(query, uid.String())
	result := new(User)
	switch err := row.Scan(&result.Username); err {
	case nil:
		result.UID = uid
		return result, nil
	case sql.ErrNoRows:
		return nil, errNotFound
	default:
		return nil, err
	}
}

func (db *db) create(username, password string) (*User, error) {
	user := new(User)

	query := "INSERT INTO users (uid, username, password_hash) VALUES ($1, $2, $3)"
	uid := uuid.New()
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user.UID = uid
	user.Username = username

	result, err := db.Exec(query, user.UID.String(), username, passwordHash)
	nRows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if nRows == 0 {
		return nil, errNotCreated
	}

	return user, nil
}

func (db *db) update(uid uuid.UUID, newPassword string) error {
	query := "UPDATE users SET password_hash=$1 WHERE uid=$2"
	passwordHash, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	result, err := db.Exec(query, passwordHash, uid.String())
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
	query := "DELETE FROM users WHERE uid=$1"
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

func (db *db) checkPassword(uid uuid.UUID, password string) (bool, error) {
	query := "SELECT password_hash FROM users WHERE uid=$1"
	row := db.QueryRow(query, uid.String())
	var passwordFromDB string
	switch err := row.Scan(&passwordFromDB); err {
	case nil:
		return checkPasswordHash(password, passwordFromDB), nil
	case sql.ErrNoRows:
		return false, errNotFound
	default:
		return false, err
	}
}

func (db *db) getUIDByUsername(username string) (uuid.UUID, error) {
	query := "SELECT uid FROM users WHERE username=$1"
	row := db.QueryRow(query, username)
	var uid string
	switch err := row.Scan(&uid); err {
	case nil:
		return uuid.Parse(uid)
	case sql.ErrNoRows:
		return uuid.Nil, errNotFound
	default:
		return uuid.Nil, err
	}
}
