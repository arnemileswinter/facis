package test

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		t.Fatalf("DATABASE_URL isn't set")
	}

	db, err := sqlx.Connect("postgres", databaseUrl)
	if err != nil {
		log.Fatalln(err)
	}

	t.Cleanup(func() { db.Close() })

	return db
}
