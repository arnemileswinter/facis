package test

import (
	"context"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
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

func createTestContractTemplate(t *testing.T, db *sqlx.DB, did *string, state templatestate.TemplateState, createdBy string) {
	name := "Test Contract Template"
	description := "Test Description"

	templateData := map[string]interface{}{
		"key": "value",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	ctx := context.Background()

	cmd := command.CreateTemplateContractCommand{
		DID:          *did,
		CreatedBy:    createdBy,
		Name:         &name,
		Description:  &description,
		TemplateData: &jsonTemplateData,
	}
	createHandler := command.CreateTemplateContractHandler{
		Ctx: ctx,
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
}

func createTestContractTemplateWithData(t *testing.T, db *sqlx.DB, did *string, state templatestate.TemplateState, createdBy string, documentNumber int, version int, name string, description string, templateData map[string]interface{}) {
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	ctx := context.Background()

	cmd := command.CreateTemplateContractCommand{
		DID:          *did,
		CreatedBy:    createdBy,
		Name:         &name,
		Description:  &description,
		TemplateData: &jsonTemplateData,
	}
	createHandler := command.CreateTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = createHandler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to create template contract: %v", err)
	}

	updateStatement := `UPDATE contract_templates SET
        	state = $2, document_number = $3, version = $4
    	WHERE did = $1
`

	_, err = db.Exec(updateStatement, *did, state, documentNumber, version)
	if err != nil {
		t.Fatalf("Failed to update template state: %v", err)
	}
}
