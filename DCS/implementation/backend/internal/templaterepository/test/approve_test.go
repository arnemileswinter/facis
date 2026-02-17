package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/templaterepository/command"
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
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, templatestate.Reviewed, db)

	ctx := context.Background()
	approvedBy := "Test User"

	cmd := command.ApproveTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		ApprovedBy:     approvedBy,
		DecisionNotes:  []string{},
	}
	handler := command.ApproveTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy := "Test User"

	qry := contracttemplate.GetContractTemplateByIdQuery{
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

	assert.Equal(t, templatestate.Approved, contractTemplate.State)
}

func TestApprove_ApproveContractTemplateInDraftState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, templatestate.Draft, db)

	ctx := context.Background()
	approvedBy := "Test User"

	cmd := command.ApproveTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		ApprovedBy:     approvedBy,
		DecisionNotes:  []string{},
	}
	handler := command.ApproveTemplateContractHandler{
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
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, templatestate.Approved, db)

	ctx := context.Background()
	approvedBy := "Test User"

	cmd := command.ApproveTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		ApprovedBy:     approvedBy,
		DecisionNotes:  []string{},
	}
	handler := command.ApproveTemplateContractHandler{
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
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, templatestate.Reviewed, db)

	ctx := context.Background()
	approvedBy := "Test User"

	cmd := command.ApproveTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now().Add(-5 * time.Second),
		ApprovedBy:     approvedBy,
		DecisionNotes:  []string{},
	}
	handler := command.ApproveTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}
