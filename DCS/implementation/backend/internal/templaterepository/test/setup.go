package test

import (
	"context"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
	"testing"

	"github.com/jmoiron/sqlx"
)

func cleanupContractTemplateTable(t *testing.T, db *sqlx.DB) {
	cleanApprovalTasksStatement := `
	-- noinspection SqlWithoutWhere
	DELETE FROM contract_templates_approval_task;
`
	_, err := db.Exec(cleanApprovalTasksStatement)
	if err != nil {
		t.Fatalf("Failed to clean table: %v", err)
	}

	cleanReviewTasksStatement := `
	-- noinspection SqlWithoutWhere
	DELETE FROM contract_templates_review_task;
`
	_, err = db.Exec(cleanReviewTasksStatement)
	if err != nil {
		t.Fatalf("Failed to clean table: %v", err)
	}

	cleanTableStatement := `
	-- noinspection SqlWithoutWhere
	DELETE FROM contract_templates;
`
	_, err = db.Exec(cleanTableStatement)
	if err != nil {
		t.Fatalf("Failed to clean table: %v", err)
	}
}

func createTestContractTemplate(t *testing.T, did *string, state templatestate.TemplateState, db *sqlx.DB) {
	name := "Test Contract Template"
	description := "Test Description"

	templateData := map[string]interface{}{
		"key": "value",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	createBy := "Test User"

	cmd := command.CreateTemplateContractCommand{
		DID:          *did,
		CreatedBy:    createBy,
		Name:         &name,
		Description:  &description,
		TemplateData: &jsonTemplateData,
	}
	createHandler := command.CreateTemplateContractHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = createHandler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to create template contract: %v", err)
	}

	updateStatement := `UPDATE contract_templates SET
        	state = $2
    	WHERE did = $1
`

	_, err = db.Exec(updateStatement, cmd.DID, state)
	if err != nil {
		t.Fatalf("Failed to update template state: %v", err)
	}

	ctx := context.Background()
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
	_, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}
}

func createTestContractTemplateWithTemplateData(t *testing.T, did *string, state templatestate.TemplateState, db *sqlx.DB, templateData map[string]interface{}) {
	name := "Test Contract Template"
	description := "Test Description"

	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	createBy := "Test User"

	cmd := command.CreateTemplateContractCommand{
		DID:          *did,
		CreatedBy:    createBy,
		Name:         &name,
		Description:  &description,
		TemplateData: &jsonTemplateData,
	}
	createHandler := command.CreateTemplateContractHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = createHandler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to create template contract: %v", err)
	}

	updateStatement := `UPDATE contract_templates SET
        	state = $2
    	WHERE did = $1
`

	_, err = db.Exec(updateStatement, cmd.DID, state)
	if err != nil {
		t.Fatalf("Failed to update template state: %v", err)
	}

	ctx := context.Background()
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
	_, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}
}
