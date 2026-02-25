package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/templaterepository/approvaltask"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/reviewtask"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVerify_VerifyContractTemplateAsReviewer(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Submitted, creator)

	ctx := context.Background()

	reviewers := []string{"Test User 1"}
	createReviewTasks(t, ctx, db, *did, reviewtaskstate.Open, creator, reviewers)

	cmd := command.VerifyCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		VerifiedBy:     reviewers[0],
		UpdatedAt:      time.Now(),
	}
	handler := command.VerifyHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to verify contract template: %v", err)
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		t.Fatal("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	exists, err := reviewtask.ExistTasksInStates(ctx, tx, *did, 1, 1, reviewtaskstate.Verified)
	if err != nil {
		t.Fatalf("Failed to check existence of review tasks: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal("could not commit transaction: %w", err)
	}

	assert.True(t, exists)
}

func TestVerify_VerifyNonExistingContractTemplate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	ctx := context.Background()

	cmd := command.VerifyCommand{
		DID:            *did,
		DocumentNumber: 2,
		Version:        2,
		UpdatedAt:      time.Now(),
		VerifiedBy:     "Test User 1",
	}
	handler := command.VerifyHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestVerify_VerifyContractTemplateAsApprover(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Submitted, creator)

	ctx := context.Background()

	approver := "Test User 1"
	createApprovalTasks(t, ctx, db, *did, approvaltaskstate.Open, creator, approver)

	cmd := command.VerifyCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		VerifiedBy:     approver,
		UpdatedAt:      time.Now(),
	}
	handler := command.VerifyHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to verify contract template: %v", err)
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		t.Fatal("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	exists, err := approvaltask.HasTaskInState(ctx, tx, *did, 1, 1, approver, approvaltaskstate.Verified)
	if err != nil {
		t.Fatalf("Failed to check existence of approval tasks: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal("could not commit transaction: %w", err)
	}

	assert.True(t, exists)
}

func TestVerify_VerifyContractTemplateAfterUpdate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Submitted, creator)

	ctx := context.Background()

	cmd := command.VerifyCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		VerifiedBy:     creator,
		UpdatedAt:      time.Now().Add(-5 * time.Second),
	}
	handler := command.VerifyHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}
