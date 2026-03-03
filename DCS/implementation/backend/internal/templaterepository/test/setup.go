package test

import (
	"context"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/command"
	approvaltask3 "digital-contracting-service/internal/templaterepository/datatype/approvaltask"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	reviewtask3 "digital-contracting-service/internal/templaterepository/datatype/reviewtask"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/datatype/templatetype"
	"digital-contracting-service/internal/templaterepository/db"
	"digital-contracting-service/internal/templaterepository/db/pg/approvaltask"
	"digital-contracting-service/internal/templaterepository/db/pg/reviewtask"
	"digital-contracting-service/internal/templaterepository/db/pg/templaterepository"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type TestRepo struct {
	CTRepo db.TemplateRepository
	RTRepo db.ReviewTaskRepo
	ATRepo db.ApprovalTaskRepo
}

func setupTestDB(t *testing.T) *sqlx.DB {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Fatalf("DATABASE_URL isn't set")
	}

	database, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		log.Fatalln(err)
	}

	t.Cleanup(func() { database.Close() })

	return database
}

func NewTestRepo(ctx context.Context) *TestRepo {
	return &TestRepo{
		CTRepo: &templaterepository.PostgresContractTemplateRepo{Ctx: ctx},
		RTRepo: &reviewtask.PostgresReviewTaskRepo{Ctx: ctx},
		ATRepo: &approvaltask.PostgresApprovalTaskRepo{Ctx: ctx},
	}
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

func createContractTemplate(t *testing.T, db *sqlx.DB, repo *TestRepo, did *string, state templatestate.TemplateState, createdBy string) {
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

	cmd := command.CreateCmd{
		DID:          *did,
		CreatedBy:    createdBy,
		TemplateType: templatetype.FrameContract,
		Name:         &name,
		Description:  &description,
		TemplateData: &jsonTemplateData,
	}
	createHandler := command.Creator{
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
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

func createTestContractTemplateWithData(t *testing.T, db *sqlx.DB, repo *TestRepo, did *string, state templatestate.TemplateState, createdBy string, documentNumber int, version int, name string, description string, templateData map[string]interface{}) {
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	ctx := context.Background()

	cmd := command.CreateCmd{
		DID:          *did,
		CreatedBy:    createdBy,
		TemplateType: templatetype.FrameContract,
		Name:         &name,
		Description:  &description,
		TemplateData: &jsonTemplateData,
	}
	createHandler := command.Creator{
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
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

func createReviewTasks(t *testing.T, ctx context.Context, db *sqlx.DB, repo *TestRepo, did string, state reviewtaskstate.ReviewTaskState, submittedBy string, reviewers []string) {
	tx, err := db.BeginTxx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	for _, reviewer := range reviewers {
		reviewTask := reviewtask3.TaskData{
			DID:            did,
			DocumentNumber: 1,
			Version:        1,
			Reviewer:       reviewer,
			State:          state,
			CreatedBy:      submittedBy,
		}
		_, err = repo.RTRepo.Create(tx, reviewTask)
		if err != nil {
			t.Fatalf("Failed to create review task: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}
}

func createApprovalTasks(t *testing.T, ctx context.Context, db *sqlx.DB, repo *TestRepo, did string, state approvaltaskstate.ApprovalTaskState, submittedBy string, approver string) {
	tx, err := db.BeginTxx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	approvalTask := approvaltask3.TaskData{
		DID:            did,
		DocumentNumber: 1,
		Version:        1,
		Approver:       approver,
		State:          state,
		CreatedBy:      submittedBy,
	}
	_, err = repo.ATRepo.Create(tx, approvalTask)
	if err != nil {
		t.Fatalf("Failed to create review task: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}
}
