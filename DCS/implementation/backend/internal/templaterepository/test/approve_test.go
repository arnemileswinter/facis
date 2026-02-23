package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApprove_ApproveContractTemplateInReviewedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Reviewed, creator)

	ctx := context.Background()

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
	defer cancel()

	approver := "Test User 1"

	createApprovalTasks(t, ctxTx, db, *did, approvaltaskstate.Open, creator, approver)

	verifyCmd := command.VerifyCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		VerifiedBy:     approver,
	}
	verifyHandler := command.VerifyHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = verifyHandler.Handle(verifyCmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	cmd := command.ApproveCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		ApprovedBy:     approver,
		DecisionNotes:  []string{},
	}
	handler := command.ApproveHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	qry := contracttemplate.GetByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
	}
	queryHandler := contracttemplate.GetByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Approved, contractTemplate.State)
}

func TestApprove_ApproveContractTemplateInReviewedStateWithoutVerifying(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Reviewed, creator)

	ctx := context.Background()

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
	defer cancel()

	approver := "Test User 1"

	createApprovalTasks(t, ctxTx, db, *did, approvaltaskstate.Open, creator, approver)

	cmd := command.ApproveCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		ApprovedBy:     approver,
		DecisionNotes:  []string{},
	}
	handler := command.ApproveHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestApprove_ApproveContractTemplateInReviewedStateWithInvalidUser(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	ctx := context.Background()

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
	defer cancel()

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Reviewed, creator)

	createApprovalTasks(t, ctxTx, db, *did, approvaltaskstate.Open, creator, "Test User 1")

	cmd := command.ApproveCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		ApprovedBy:     "Test User 2",
		DecisionNotes:  []string{},
	}
	handler := command.ApproveHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.Error(t, err)
}

func TestApprove_ApproveContractTemplateInDraftState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Draft, creator)

	ctx := context.Background()

	cmd := command.ApproveCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		ApprovedBy:     "Test User 1",
		DecisionNotes:  []string{},
	}
	handler := command.ApproveHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestApprove_ApproveContractTemplateInApprovedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Approved, creator)

	ctx := context.Background()

	cmd := command.ApproveCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		ApprovedBy:     "Test User 1",
		DecisionNotes:  []string{},
	}
	handler := command.ApproveHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestApprove_ApproveContractTemplateAfterUpdate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Reviewed, creator)

	ctx := context.Background()

	cmd := command.ApproveCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now().Add(-5 * time.Second),
		ApprovedBy:     "Test User 1",
		DecisionNotes:  []string{},
	}
	handler := command.ApproveHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}
