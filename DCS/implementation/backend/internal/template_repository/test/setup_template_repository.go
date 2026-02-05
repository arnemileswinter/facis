package test

import (
	"context"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/command"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"digital-contracting-service/internal/template_repository/query"
	"testing"

	"github.com/jmoiron/sqlx"
)

func dropAndCreateContractTemplateTable(t *testing.T, db *sqlx.DB) {
	dropStatement := `DROP TABLE IF EXISTS contract_templates`

	createStatement := `
    CREATE TABLE IF NOT EXISTS contract_templates (
    did VARCHAR(255) PRIMARY KEY CHECK (did <> '' AND did IS NOT NULL),
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    document_number INT NOT NULL,
    version INT NOT NULL,
    state VARCHAR(16) NOT NULL,

    name VARCHAR(255) NOT NULL,
    description TEXT,
    meta_data JSONB
);
`

	_, err := db.Exec(dropStatement)
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}

	_, err = db.Exec(createStatement)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
}

func createTestContractTemplate(t *testing.T, did *string, state template_state.TemplateState, db *sqlx.DB) {
	name := "Test Contract Template"
	description := "Test Description"

	metaData := map[string]interface{}{
		"key": "value",
	}
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
		Db: db,
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

	qry := query.GetContractTemplateByIdQuery{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	queryHandler := query.GetContractTemplateByIdHandler{
		Ctx: ctx,
		Db:  db,
	}
	_, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}
}
