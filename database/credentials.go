package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raian621/dump/models/storage"
)

func InsertCredentials(db *pgxpool.Pool, creds *storage.Credentials) error {
	_, err := db.Exec(
		context.Background(),
		"INSERT INTO credentials (user_id, username, passhash) VALUES ($1, $2, $3)",
		creds.UserId, creds.Username, creds.Passhash)
	return err
}

func GetPasshashForUsername(db *pgxpool.Pool, username string) (string, error) {
	var passhash string
	row := db.QueryRow(
		context.Background(),
		"SELECT passhash FROM credentials WHERE username = $1", username)
	if err := row.Scan(&passhash); err != nil {
		return "", err
	}
	return passhash, nil
}
