package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/template_repository/command"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"digital-contracting-service/internal/template_repository/query"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubmit_RejectContractTemplateInReviewedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Reviewed, db)

	ctx := context.Background()
	rejectedBy := "Test User"

	cmd := command.RejectTemplateContractCommand{
		DID:        *did,
		RejectedBy: rejectedBy,
		Reason:     "Test Reason",
	}
	handler := command.RejectTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy := "Test User"

	qry := query.GetContractTemplateByIdQuery{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	queryHandler := query.GetContractTemplateByIdHandler{
		Ctx: ctx,
		Db:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, template_state.Draft, contractTemplate.State)
}

func TestSubmit_RejectContractTemplateInDraftState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Draft, db)

	ctx := context.Background()
	rejectedBy := "Test User"

	cmd := command.RejectTemplateContractCommand{
		DID:        *did,
		RejectedBy: rejectedBy,
		Reason:     "Test Reason",
	}
	handler := command.RejectTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestSubmit_RejectContractTemplateInApprovedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Approved, db)

	ctx := context.Background()
	rejectedBy := "Test User"

	cmd := command.RejectTemplateContractCommand{
		DID:        *did,
		RejectedBy: rejectedBy,
		Reason:     "Test Reason",
	}
	handler := command.RejectTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}
