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

func TestRegister_RegisterContractTemplateDataInValidState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Approved, creator)

	ctx := context.Background()

	cmd := command.RegisterCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.RegisterHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	qry := contracttemplate.GetByIDQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
	}
	queryHandler := contracttemplate.GetByIDHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query contract template: %v", err)
	}

	assert.Equal(t, contractTemplate.DID, *did)
	assert.Equal(t, templatestate.Registered, contractTemplate.State)
}

func TestRegister_RegisterNonExistingContractTemplate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	ctx := context.Background()

	cmd := command.RegisterCommand{
		DID:            *did,
		DocumentNumber: 2,
		Version:        2,
		UpdatedAt:      time.Now(),
		RegisteredBy:   "Test User 1",
	}
	handler := command.RegisterHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestRegister_RegisterContractTemplateDataInDraftState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Draft, creator)

	ctx := context.Background()

	cmd := command.RegisterCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.RegisterHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestRegister_RegisterContractTemplateDataInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Submitted, creator)

	ctx := context.Background()

	cmd := command.RegisterCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.RegisterHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestRegister_RegisterContractTemplateDataInRejectedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Rejected, creator)

	ctx := context.Background()

	cmd := command.RegisterCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.RegisterHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestRegister_RegisterContractTemplateDataInReviewedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Reviewed, creator)

	ctx := context.Background()

	cmd := command.RegisterCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.RegisterHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestRegister_RegisterContractTemplateDataInRegisteredState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Registered, creator)

	ctx := context.Background()

	cmd := command.RegisterCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.RegisterHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestRegister_RegisterContractTemplateDataInArchivedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Archived, creator)

	ctx := context.Background()

	cmd := command.RegisterCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.RegisterHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}
