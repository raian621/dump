package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raian621/dump/models/storage"
)

func InsertUser(db *pgxpool.Pool, user *storage.User) {
	db.Exec(
		context.Background(), "INSERT INTO users (id, username) VALUES($1, $2)",
		user.Id, user.Username)
}

func GetUserById(db *pgxpool.Pool, id uint64) (*storage.User, error) {
	user := &storage.User{}
	row := db.QueryRow(context.Background(), "SELECT id, username FROM users")
	if err := row.Scan(&user.Id, &user.Username); err != nil {
		return nil, err
	}
	return user, nil
}

func UsernameExists(db *pgxpool.Pool, username string) (bool, error) {
	// consider using a bloom filter if this becomes a bottleneck
	var exists bool
	row := db.QueryRow(context.Background(),
		"SELECT COUNT(*) > 0 FROM users WHERE username = $1", username)
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func CreateUserWithCredentials(db *pgxpool.Pool, creds *storage.Credentials) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(),
		"INSERT INTO users (username) VALUES ($1)",
		creds.Username)
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO credentials (user_id, username, passhash) VALUES ((SELECT id FROM users WHERE username = $1), $1, $2)",
		creds.Username, creds.Passhash)
	if err != nil {
		return err
	}
	return tx.Commit(context.Background())
}

func GetUserIdFromUsername(db *pgxpool.Pool, username string) (userId int32, err error) {
	row := db.QueryRow(
		context.Background(), "SELECT id FROM users WHERE username = $1", username)
	err = row.Scan(&userId)
	return userId, err
}
