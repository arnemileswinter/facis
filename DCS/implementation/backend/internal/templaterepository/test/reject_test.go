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

func TestCreate_RejectContractTemplateInReviewedState(t *testing.T) {

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

	cmd := command.RejectCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		RejectedBy:     approver,
		Reason:         "Test Reason",
	}
	handler := command.RejectHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	retrievedBy := "Test User"

	qry := contracttemplate.GetByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query contract template: %v", err)
	}

	assert.Equal(t, templatestate.Draft, contractTemplate.State)
}

func TestCreate_RejectContractTemplateInReviewedStateWithInvalidUser(t *testing.T) {

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

	createApprovalTasks(t, ctxTx, db, *did, approvaltaskstate.Open, creator, "Test User 1")

	cmd := command.RejectCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		RejectedBy:     "Test User 2",
		Reason:         "Test Reason",
	}
	handler := command.RejectHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestCreate_RejectNonExistingContractTemplate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	ctx := context.Background()

	cmd := command.RejectCommand{
		DID:            *did,
		DocumentNumber: 2,
		Version:        2,
		UpdatedAt:      time.Now(),
		RejectedBy:     "Test User 1",
	}
	handler := command.RejectHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestCreate_RejectContractTemplateInDraftState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	createContractTemplate(t, db, did, templatestate.Draft, "Test User")

	ctx := context.Background()
	rejectedBy := "Test User"

	cmd := command.RejectCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		RejectedBy:     rejectedBy,
		Reason:         "Test Reason",
	}
	handler := command.RejectHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestCreate_RejectContractTemplateInApprovedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	createContractTemplate(t, db, did, templatestate.Approved, "Test User")

	ctx := context.Background()
	rejectedBy := "Test User"

	cmd := command.RejectCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		RejectedBy:     rejectedBy,
		Reason:         "Test Reason",
	}
	handler := command.RejectHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestCreate_RejectContractTemplateAfterUpdate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	createContractTemplate(t, db, did, templatestate.Reviewed, "Test User")

	ctx := context.Background()
	rejectedBy := "Test User"

	cmd := command.RejectCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now().Add(-5 * time.Second),
		RejectedBy:     rejectedBy,
		Reason:         "Test Reason",
	}
	handler := command.RejectHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}
