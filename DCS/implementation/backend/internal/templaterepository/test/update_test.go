package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdate_UpdateContractTemplateDataInDraftState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Draft, creator)

	ctx := context.Background()
	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateTemplateContractHandler{
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

	assert.Equal(t, *did, contractTemplate.DID)
	assert.Equal(t, name, *contractTemplate.Name)
	assert.Equal(t, description, *contractTemplate.Description)
	//assert.Equal(t, jsonTemplateData, contractTemplate.TemplateData)
}

func TestUpdate_UpdateContractTemplateDataInDraftStateWithInvalidUser(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Draft, creator)

	ctx := context.Background()
	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      "Test User 1",
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateDataInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Submitted, creator)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateTemplateContractHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateDataInSubmittedStateWithInvalidUser(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Submitted, creator)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      "Test User 1",
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateTemplateContractHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateDataInDraftApprovedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Approved, creator)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateTemplateContractHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateDataInDraftApprovedStateWithInvalidUser(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Approved, creator)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      "Test User 1",
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateTemplateContractHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateAfterUpdate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Draft, creator)

	ctx := context.Background()
	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now().Add(-5 * time.Second),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}
