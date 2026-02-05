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

	dropAndCreateContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	currentContractState := template_state.Draft
	createTestContractTemplate(t, did, currentContractState, db)

	submittedBy := "Test User"

	cmd := command.SubmitTemplateContractCommand{
		DID:                          *did,
		SubmittedBy:                  submittedBy,
		CurrentContractTemplateState: currentContractState,
		ActionFlag:                   nil,
		ReviewComments:               []string{},
	}
	handler := command.SubmitTemplateContractHandler{
		Db: db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	ctx := context.Background()
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

	assert.Equal(t, *did, contractTemplate.DID)
	assert.Equal(t, template_state.Submitted, contractTemplate.State)
}

func TestSubmit_ApproveContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	dropAndCreateContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	currentContractState := template_state.Submitted
	createTestContractTemplate(t, did, currentContractState, db)

	submittedBy := "Test User"
	actionFlag := action_flag.Approval

	cmd := command.SubmitTemplateContractCommand{
		DID:                          *did,
		SubmittedBy:                  submittedBy,
		CurrentContractTemplateState: currentContractState,
		ActionFlag:                   &actionFlag,
		ReviewComments:               []string{},
	}
	handler := command.SubmitTemplateContractHandler{
		Db: db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	ctx := context.Background()
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

	assert.Equal(t, *did, contractTemplate.DID)
	assert.Equal(t, template_state.Reviewed, contractTemplate.State)
}

func TestSubmit_DeclineContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	dropAndCreateContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	currentContractState := template_state.Submitted
	createTestContractTemplate(t, did, currentContractState, db)

	submittedBy := "Test User"
	actionFlag := action_flag.Draft

	cmd := command.SubmitTemplateContractCommand{
		DID:                          *did,
		SubmittedBy:                  submittedBy,
		CurrentContractTemplateState: currentContractState,
		ActionFlag:                   &actionFlag,
		ReviewComments:               []string{},
	}
	handler := command.SubmitTemplateContractHandler{
		Db: db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	ctx := context.Background()
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

	assert.Equal(t, *did, contractTemplate.DID)
	assert.Equal(t, template_state.Draft, contractTemplate.State)
}

func TestSubmit_SubmitContractTemplateInSubmittedStateWithoutActionFlag(t *testing.T) {

	db := setupTestDB(t)

	dropAndCreateContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	currentContractState := template_state.Submitted
	createTestContractTemplate(t, did, currentContractState, db)

	submittedBy := "Test User"

	cmd := command.SubmitTemplateContractCommand{
		DID:                          *did,
		SubmittedBy:                  submittedBy,
		CurrentContractTemplateState: currentContractState,
		ActionFlag:                   nil,
		ReviewComments:               []string{},
	}
	handler := command.SubmitTemplateContractHandler{
		Db: db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestSubmit_DeclineContractTemplateInReviewedState(t *testing.T) {

	db := setupTestDB(t)

	dropAndCreateContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	currentContractState := template_state.Reviewed
	createTestContractTemplate(t, did, currentContractState, db)

	submittedBy := "Test User"
	actionFlag := action_flag.Draft

	cmd := command.SubmitTemplateContractCommand{
		DID:                          *did,
		SubmittedBy:                  submittedBy,
		CurrentContractTemplateState: currentContractState,
		ActionFlag:                   &actionFlag,
		ReviewComments:               []string{},
	}
	handler := command.SubmitTemplateContractHandler{
		Db: db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	ctx := context.Background()
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

	assert.Equal(t, *did, contractTemplate.DID)
	assert.Equal(t, template_state.Draft, contractTemplate.State)
}

func TestSubmit_ApproveContractTemplateInReviewedState(t *testing.T) {

	db := setupTestDB(t)

	dropAndCreateContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	currentContractState := template_state.Reviewed
	createTestContractTemplate(t, did, currentContractState, db)

	submittedBy := "Test User"
	actionFlag := action_flag.Approval

	cmd := command.SubmitTemplateContractCommand{
		DID:                          *did,
		SubmittedBy:                  submittedBy,
		CurrentContractTemplateState: currentContractState,
		ActionFlag:                   &actionFlag,
		ReviewComments:               []string{},
	}
	handler := command.SubmitTemplateContractHandler{
		Db: db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestSubmit_SubmitContractTemplateInApprovedStateWithoutActionFlag(t *testing.T) {

	db := setupTestDB(t)

	dropAndCreateContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	currentContractState := template_state.Approved
	createTestContractTemplate(t, did, currentContractState, db)

	submittedBy := "Test User"

	cmd := command.SubmitTemplateContractCommand{
		DID:                          *did,
		SubmittedBy:                  submittedBy,
		CurrentContractTemplateState: currentContractState,
		ActionFlag:                   nil,
		ReviewComments:               []string{},
	}
	handler := command.SubmitTemplateContractHandler{
		Db: db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}
