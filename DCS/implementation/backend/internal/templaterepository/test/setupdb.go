package test

import (
	"context"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
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

func createReviewTasks(t *testing.T, ctx context.Context, db *sqlx.DB, did string, state reviewtaskstate.ReviewTaskState, submittedBy string, reviewers []string) {
	tx, err := db.BeginTxx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	for _, reviewer := range reviewers {
		reviewTask := templaterepository.ReviewTaskData{
			DID:            did,
			DocumentNumber: 1,
			Version:        1,
			Reviewer:       reviewer,
			State:          state,
			CreatedBy:      submittedBy,
		}
		_, err = templaterepository.CreateReviewTasks(ctx, tx, reviewTask)
		if err != nil {
			t.Fatalf("Failed to create review task: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}
}

func createApprovalTasks(t *testing.T, ctx context.Context, db *sqlx.DB, did string, state approvaltaskstate.ApprovalTaskState, submittedBy string, approver string) {
	tx, err := db.BeginTxx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	approvalTask := templaterepository.ApprovalTaskData{
		DID:            did,
		DocumentNumber: 1,
		Version:        1,
		Approver:       approver,
		State:          state,
		CreatedBy:      submittedBy,
	}
	_, err = templaterepository.CreateApprovalTask(ctx, tx, approvalTask)
	if err != nil {
		t.Fatalf("Failed to create review task: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}
}
