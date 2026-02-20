package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
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

	templateData := map[string]interface{}{}
	jsonMetaData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	ctx := context.Background()

	creator := "Test User"

	cmd := command.CreateTemplateContractCommand{
		DID:          *did,
		CreatedBy:    creator,
		Name:         &name,
		Description:  &description,
		TemplateData: &jsonMetaData,
	}
	createHandler := command.CreateTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = createHandler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to create template contract: %v", err)
	}

	qry := contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
	assert.Equal(t, jsonMetaData, *contractTemplate.TemplateData)
}
