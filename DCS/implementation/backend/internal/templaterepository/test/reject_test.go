package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
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

	createTestContractTemplate(t, did, templatestate.Reviewed, db)

	ctx := context.Background()
	rejectedBy := "Test User"

	cmd := command.RejectTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RejectedBy:     rejectedBy,
		Reason:         "Test Reason",
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

	qry := contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Draft, contractTemplate.State)
}

func TestSubmit_RejectContractTemplateInDraftState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, templatestate.Draft, db)

	ctx := context.Background()
	rejectedBy := "Test User"

	cmd := command.RejectTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RejectedBy:     rejectedBy,
		Reason:         "Test Reason",
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

	createTestContractTemplate(t, did, templatestate.Approved, db)

	ctx := context.Background()
	rejectedBy := "Test User"

	cmd := command.RejectTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RejectedBy:     rejectedBy,
		Reason:         "Test Reason",
	}
	handler := command.RejectTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}
