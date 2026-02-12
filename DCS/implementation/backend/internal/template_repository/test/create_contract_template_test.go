package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/command"
	"digital-contracting-service/internal/template_repository/query"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate_CreateNewContractTemplate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	name := "Test Contract Template"
	description := "Test Description"

	metaData := map[string]interface{}{}
	jsonMetaData, err := datatype.NewJSON(metaData)
	if err != nil {
		t.Fatalf("Failed to create JSON metadata: %v", err)
	}

	createBy := "Test User"

	cmd := command.CreateTemplateContractCommand{
		DID:         *did,
		CreatedBy:   createBy,
		Name:        &name,
		Description: &description,
		MetaData:    &jsonMetaData,
	}
	createHandler := command.CreateTemplateContractHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = createHandler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to create template contract: %v", err)
	}

	ctx := context.Background()
	retrievedBy := "Test User"

	qry := query.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := query.GetContractTemplateByIdHandler{
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
	assert.Equal(t, jsonMetaData, *contractTemplate.MetaData)
}
