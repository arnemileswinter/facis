package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/command"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"digital-contracting-service/internal/template_repository/query"
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

	currentContractState := template_state.Draft
	createTestContractTemplate(t, did, currentContractState, db)

	updateBy := "Test User"
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
		UpdatedBy:   updateBy,
		Name:        &name,
		Description: &description,
		MetaData:    &jsonMetaData,
	}
	handler := command.UpdateTemplateContractHandler{
		Ctx: context.Background(),
		DB:  db,
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

	assert.Equal(t, name, contractTemplate.Name)
	assert.Equal(t, description, contractTemplate.Description)
	//assert.Equal(t, jsonMetaData, contractTemplate.MetaData)
}

func TestSubmit_UpdateContractTemplateMetaDataInDraftSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	currentContractState := template_state.Submitted
	createTestContractTemplate(t, did, currentContractState, db)

	updateBy := "Test User"
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
		UpdatedBy:   updateBy,
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

	currentContractState := template_state.Approved
	createTestContractTemplate(t, did, currentContractState, db)

	updateBy := "Test User"
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
		UpdatedBy:   updateBy,
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
