package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/template_repository/command"
	"digital-contracting-service/internal/template_repository/datatype/action_flag"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"digital-contracting-service/internal/template_repository/query"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubmit_SubmitContractTemplateInDraftState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Draft, db)

	ctx := context.Background()
	submittedBy := "Test User"

	cmd := command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    submittedBy,
		ActionFlag:     nil,
		ReviewComments: []string{},
	}
	handler := command.SubmitTemplateContractHandler{
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

	assert.Equal(t, template_state.Submitted, contractTemplate.State)
}

func TestSubmit_ApproveContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Submitted, db)

	ctx := context.Background()
	submittedBy := "Test User"
	actionFlag := action_flag.Approval

	cmd := command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    submittedBy,
		ActionFlag:     &actionFlag,
		ReviewComments: []string{},
	}
	handler := command.SubmitTemplateContractHandler{
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

	assert.Equal(t, template_state.Reviewed, contractTemplate.State)
}

func TestSubmit_DeclineContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Submitted, db)

	ctx := context.Background()
	submittedBy := "Test User"
	actionFlag := action_flag.Draft

	cmd := command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    submittedBy,
		ActionFlag:     &actionFlag,
		ReviewComments: []string{},
	}
	handler := command.SubmitTemplateContractHandler{
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

func TestSubmit_SubmitContractTemplateInSubmittedStateWithoutActionFlag(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Submitted, db)

	ctx := context.Background()
	submittedBy := "Test User"

	cmd := command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    submittedBy,
		ActionFlag:     nil,
		ReviewComments: []string{},
	}
	handler := command.SubmitTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestSubmit_SubmitContractTemplateInReviewedStateToResubmission(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Reviewed, db)

	ctx := context.Background()
	submittedBy := "Test User"

	cmd := command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    submittedBy,
		ActionFlag:     nil,
		ReviewComments: []string{},
	}
	handler := command.SubmitTemplateContractHandler{
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

	assert.Equal(t, template_state.Submitted, contractTemplate.State)
}
