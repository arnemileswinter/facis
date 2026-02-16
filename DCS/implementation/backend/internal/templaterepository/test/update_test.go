package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubmit_UpdateContractTemplateMetaDataInDraftState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	currentContractState := templatestate.Draft
	createTestContractTemplate(t, did, currentContractState, db)

	ctx := context.Background()
	metaData := map[string]interface{}{
		"test": "update",
	}
	jsonMetaData, err := datatype.NewJSON(metaData)
	if err != nil {
		t.Fatalf("Failed to create JSON metadata: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		Name:           &name,
		Description:    &description,
		MetaData:       &jsonMetaData,
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

	assert.Equal(t, *did, contractTemplate.DID)
	assert.Equal(t, name, *contractTemplate.Name)
	assert.Equal(t, description, *contractTemplate.Description)
	//assert.Equal(t, jsonMetaData, contractTemplate.MetaData)
}

func TestSubmit_UpdateContractTemplateMetaDataInDraftSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	currentContractState := templatestate.Submitted
	createTestContractTemplate(t, did, currentContractState, db)

	metaData := map[string]interface{}{
		"test": "update",
	}
	jsonMetaData, err := datatype.NewJSON(metaData)
	if err != nil {
		t.Fatalf("Failed to create JSON metadata: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateTemplateContractCommand{
		DID:         *did,
		Name:        &name,
		Description: &description,
		MetaData:    &jsonMetaData,
	}
	handler := command.UpdateTemplateContractHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestSubmit_UpdateContractTemplateMetaDataInDraftApprovedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	currentContractState := templatestate.Approved
	createTestContractTemplate(t, did, currentContractState, db)

	metaData := map[string]interface{}{
		"test": "update",
	}
	jsonMetaData, err := datatype.NewJSON(metaData)
	if err != nil {
		t.Fatalf("Failed to create JSON metadata: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateTemplateContractCommand{
		DID:         *did,
		Name:        &name,
		Description: &description,
		MetaData:    &jsonMetaData,
	}
	handler := command.UpdateTemplateContractHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}
