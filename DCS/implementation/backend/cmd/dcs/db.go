package main

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"log"
	"os"
)

func NewDatabaseConnection(ctx context.Context) (*sqlx.DB, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		return nil, errors.New("DATABASE_URL isn't set")
	}

	log.Printf("Connecting to database")

	db, err := sqlx.Connect("postgres", databaseUrl)
	if err != nil {
		log.Fatalln(err)
	}

	return db, nil
}
