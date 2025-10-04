package database

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"path"

	"github.com/jackc/pgx/v5/pgxpool"
)

const migrationDir = "migrations"

func ApplyMigrations(db *pgxpool.Pool, migrations []string) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		panic(err)
	}
	for _, migration := range migrations {
		log.Printf("Applying migration `%s`...\n", migration)
		bytes, err := os.ReadFile(path.Join(migrationDir, migration))
		if err != nil {
			panic(err)
		}
		script := string(bytes)
		if _, err := tx.Exec(context.Background(), script); err != nil {
			panic(err)
		}
		if _, err := tx.Exec(context.Background(), "INSERT INTO migrations (migration) VALUES ($1)", migration); err != nil {
			panic(err)
		}
	}
	if err := tx.Commit(context.Background()); err != nil {
		panic(err)
	}
}

func ReadMigrationRecordsFromFile(records io.Reader) []string {
	scanner := bufio.NewScanner(records)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func ReadMigrationRecordsFromDb(db *pgxpool.Pool) ([]string, error) {
	var (
		migrations []string
		migration  string
	)

	rows, err := db.Query(context.Background(), "SELECT migration FROM migrations ORDER BY id")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(&migration); err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}

	return migrations, nil
}

func GetNeededMigrations(migrationRecords []string, dbMigrations []string) []string {
	var (
		migrations       map[string]bool = make(map[string]bool, 0)
		neededMigrations []string
	)

	for _, migration := range dbMigrations {
		migrations[migration] = true
	}
	for _, migration := range migrationRecords {
		if _, ok := migrations[migration]; ok {
			log.Printf("Skipping migration `%s`.", migration)
			continue
		}
		neededMigrations = append(neededMigrations, migration)
	}

	return neededMigrations
}
