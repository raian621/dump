package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/raian621/dump/auth"
	"github.com/raian621/dump/database"
	"github.com/raian621/dump/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	s := server.New()
	accessTtl, refreshTtl := getTokenTtls()
	s.AddTokenFactory(auth.NewTokenFactory(accessTtl, refreshTtl, getJwtSecret()))
	s.AddHandlers()
	db := getDbClient()
	s.AddDatabaseClient(db)
	applyMigrations(db)
	log.Fatalln(s.Start(":1234"))
}

func applyMigrations(db *pgxpool.Pool) {
	log.Println("Applying database migrations...")
	dbMigrations, err := database.ReadMigrationRecordsFromDb(db)
	if err != nil {
		panic(err)
	}
	records, err := os.Open("migrations/_records.txt")
	if err != nil {
		panic(err)
	}
	recordMigrations := database.ReadMigrationRecordsFromFile(records)

	neededMigrations := database.GetNeededMigrations(recordMigrations, dbMigrations)
	database.ApplyMigrations(db, neededMigrations)
}

func getDbClient() *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), getConnString())
	if err != nil {
		panic(err)
	}
	return db
}

func getConnString() string {
	return fmt.Sprintf(
		"postgresql://%s:%s/%s?user=%s&password=%s", os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"))
}

func getTokenTtls() (int, int) {
	var (
		accessTtl  = 15 * 60       // 15 minute TTL
		refreshTtl = 7 * 24 * 3600 // 7 day TTL
		err        error
	)

	if envAccessTtl, found := os.LookupEnv("JWT_ACCESS_TTL"); found {
		accessTtl, err = strconv.Atoi(envAccessTtl)
		if err != nil {
			panic(err)
		}
	}
	if envRefreshTtl, found := os.LookupEnv("JWT_REFRESH_TTL"); found {
		refreshTtl, err = strconv.Atoi(envRefreshTtl)
		if err != nil {
			panic(err)
		}
	}

	return accessTtl, refreshTtl
}

func getJwtSecret() []byte {
	secretBase64 := os.Getenv("JWT_SECRET")
	secret := make([]byte, base64.RawURLEncoding.DecodedLen(len(secretBase64)))
	n, err := base64.RawURLEncoding.Decode(secret, []byte(secretBase64))
	if err != nil {
		panic(err)
	}
	return secret[:n]
}
