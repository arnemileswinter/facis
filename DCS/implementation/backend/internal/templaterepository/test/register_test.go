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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Approved, creator)

	cmd := command.RegisterCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.Registrar{
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		ATRepo: repo.ATRepo,
		RTRepo: repo.RTRepo,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	qry := contracttemplate.GetByIDQry{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
	}
	queryHandler := contracttemplate.GetByIDHandler{
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	cmd := command.RegisterCmd{
		DID:            *did,
		DocumentNumber: 2,
		Version:        2,
		UpdatedAt:      time.Now(),
		RegisteredBy:   "Test User 1",
	}
	handler := command.Registrar{
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		ATRepo: repo.ATRepo,
		RTRepo: repo.RTRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Draft, creator)

	cmd := command.RegisterCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.Registrar{
		Ctx:    ctx,
		DB:     db,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
		CTRepo: repo.CTRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Submitted, creator)

	cmd := command.RegisterCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.Registrar{
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		ATRepo: repo.ATRepo,
		RTRepo: repo.RTRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Rejected, creator)

	cmd := command.RegisterCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.Registrar{
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		ATRepo: repo.ATRepo,
		RTRepo: repo.RTRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Reviewed, creator)

	cmd := command.RegisterCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.Registrar{
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		ATRepo: repo.ATRepo,
		RTRepo: repo.RTRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Registered, creator)

	cmd := command.RegisterCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.Registrar{
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		ATRepo: repo.ATRepo,
		RTRepo: repo.RTRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Archived, creator)

	cmd := command.RegisterCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RegisteredBy:   creator,
		UpdatedAt:      time.Now(),
	}
	handler := command.Registrar{
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		ATRepo: repo.ATRepo,
		RTRepo: repo.RTRepo,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}
