package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/reviewtask"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReview_CreateReviewTasks(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Submitted, creator)

	ctx := context.Background()

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
	defer cancel()

	assignees := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	tx, err := db.BeginTxx(ctxTx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	for _, assignee := range assignees {
		reviewTask := reviewtask.TaskData{
			DID:            *did,
			DocumentNumber: 1,
			Version:        1,
			Reviewer:       assignee,
			State:          reviewtaskstate.Open,
			CreatedBy:      creator,
		}
		_, err = reviewtask.CreateTask(ctx, tx, reviewTask)
		if err != nil {
			t.Fatalf("Failed to create review task: %v", err)
		}
	}

	exists, err := reviewtask.ExistTasksInState(ctx, tx, *did, 1, 1, reviewtaskstate.Open)
	if err != nil {
		t.Fatalf("Failed to check if review task exists: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	assert.True(t, exists)
}

func TestReview_CreateReviewTasksAndApproveThem(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Submitted, creator)

	ctx := context.Background()

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
	defer cancel()

	assignees := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	tx, err := db.BeginTxx(ctxTx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	for _, assignee := range assignees {
		reviewTask := reviewtask.TaskData{
			DID:            *did,
			DocumentNumber: 1,
			Version:        1,
			Reviewer:       assignee,
			State:          reviewtaskstate.Open,
			CreatedBy:      creator,
		}
		_, err = reviewtask.CreateTask(ctx, tx, reviewTask)
		if err != nil {
			t.Fatalf("Failed to create review task: %v", err)
		}
	}

	for _, assignee := range assignees {
		err := reviewtask.UpdateTask(ctx, tx, *did, 1, 1, assignee, reviewtaskstate.Approved)
		if err != nil {
			t.Fatalf("Failed to approve review task: %v", err)
		}
	}

	exists, err := reviewtask.ExistTasksInState(ctx, tx, *did, 1, 1, reviewtaskstate.Open)
	if err != nil {
		t.Fatalf("Failed to check if review task exists: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	assert.False(t, exists)
}
