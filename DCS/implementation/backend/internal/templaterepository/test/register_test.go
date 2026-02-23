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

	assert.Equal(t, contractTemplate.DID, *did)
	assert.Equal(t, templatestate.Registered, contractTemplate.State)
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
