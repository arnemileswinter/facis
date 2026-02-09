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

func cleanupContractTemplateTable(t *testing.T, db *sqlx.DB) {
	cleanTableStatement := "DELETE FROM contract_templates;"
	_, err := db.Exec(cleanTableStatement)
	if err != nil {
		t.Fatalf("Failed to clean table: %v", err)
	}

	//cleanTableStatement = "DELETE FROM outbox_events;"
	//_, err = db.Exec(cleanTableStatement)
	//if err != nil {
	//	t.Fatalf("Failed to clean table: %v", err)
	//}
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
