package test

import (
	"context"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/approvaltask"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/datatype/templatetype"
	"digital-contracting-service/internal/templaterepository/reviewtask"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		t.Fatalf("DATABASE_URL isn't set")
	}

	db, err := sqlx.Connect("postgres", databaseUrl)
	if err != nil {
		log.Fatalln(err)
	}

	t.Cleanup(func() { db.Close() })

	return db
}

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

func createContractTemplate(t *testing.T, db *sqlx.DB, did *string, state templatestate.TemplateState, createdBy string) {
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

	cmd := command.CreateCommand{
		DID:          *did,
		CreatedBy:    createdBy,
		TemplateType: templatetype.FrameContract,
		Name:         &name,
		Description:  &description,
		TemplateData: &jsonTemplateData,
	}
	createHandler := command.CreateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = createHandler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to create contract template: %v", err)
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

	cmd := command.CreateCommand{
		DID:          *did,
		CreatedBy:    createdBy,
		TemplateType: templatetype.FrameContract,
		Name:         &name,
		Description:  &description,
		TemplateData: &jsonTemplateData,
	}
	createHandler := command.CreateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = createHandler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to create contract template: %v", err)
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

func createReviewTasks(t *testing.T, ctx context.Context, db *sqlx.DB, did string, state reviewtaskstate.ReviewTaskState, submittedBy string, reviewers []string) {
	tx, err := db.BeginTxx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	for _, reviewer := range reviewers {
		reviewTask := reviewtask.TaskData{
			DID:            did,
			DocumentNumber: 1,
			Version:        1,
			Reviewer:       reviewer,
			State:          state,
			CreatedBy:      submittedBy,
		}
		_, err = reviewtask.CreateTask(ctx, tx, reviewTask)
		if err != nil {
			t.Fatalf("Failed to create review task: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}
}

func createApprovalTasks(t *testing.T, ctx context.Context, db *sqlx.DB, did string, state approvaltaskstate.ApprovalTaskState, submittedBy string, approver string) {
	tx, err := db.BeginTxx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	approvalTask := approvaltask.TaskData{
		DID:            did,
		DocumentNumber: 1,
		Version:        1,
		Approver:       approver,
		State:          state,
		CreatedBy:      submittedBy,
	}
	_, err = approvaltask.CreateTask(ctx, tx, approvalTask)
	if err != nil {
		t.Fatalf("Failed to create review task: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}
}
