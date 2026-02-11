package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/template_repository"
	"digital-contracting-service/internal/template_repository/command"
	"digital-contracting-service/internal/template_repository/datatype/action_flag"
	"digital-contracting-service/internal/template_repository/datatype/review_task_state"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"digital-contracting-service/internal/template_repository/query"
	"slices"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSubmit_SubmitContractTemplateInDraftState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Draft, db)

	ctx := context.Background()
	submittedBy := "Test User"

	cmd := command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    submittedBy,
		ActionFlag:     nil,
		ReviewComments: nil,
		Assignees: []string{
			"Test User 2",
			"Test User 3",
			"Test User 4",
		},
	}
	handler := command.SubmitTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy := "Test User"

	qry := query.GetContractTemplateByIdQuery{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	queryHandler := query.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, template_state.Submitted, contractTemplate.State)

	queryReviewTasks := query.GetAllContractTemplateReviewTasksForDID{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	handlerReviewTasks := query.GetAllContractTemplateReviewTasksForDIDHandler{
		Ctx: ctx,
		DB:  db,
	}
	reviewTasks, err := handlerReviewTasks.Handle(queryReviewTasks)
	if err != nil {
		t.Fatalf("Failed to query template review tasks: %v", err)
	}

	for _, reviewTask := range reviewTasks {
		assert.Equal(t, review_task_state.Open, reviewTask.State)

		if !slices.Contains(cmd.Assignees, reviewTask.Assignee) {
			t.Fatalf("Assignee not found in review tasks: %v", reviewTask)
		}
	}
}

func createReviewTasks(t *testing.T, ctx context.Context, db *sqlx.DB, did string, submittedBy string, assignees []string) error {
	tx, err := db.BeginTxx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	for _, assignee := range assignees {
		reviewTask := template_repository.ReviewTaskData{
			DID:            did,
			DocumentNumber: 1,
			Version:        1,
			Assignee:       assignee,
			State:          review_task_state.Open,
			CreatedBy:      submittedBy,
		}
		_, err = template_repository.CreateReviewTask(ctx, tx, reviewTask)
		if err != nil {
			t.Fatalf("Failed to create review task: %v", err)
		}
	}
	return tx.Commit()
}

func TestSubmit_OneAssigneeApprovedContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Submitted, db)

	ctx := context.Background()
	submittedBy := "Test User"

	assignees := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	ctxTx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = createReviewTasks(t, ctxTx, db, *did, submittedBy, assignees)
	if err != nil {
		t.Fatalf("Failed to create review tasks: %v", err)
	}

	actionFlag := action_flag.Approval

	submittedBy = "Test User 1"
	cmd := command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    submittedBy,
		ActionFlag:     &actionFlag,
		ReviewComments: []string{},
	}
	handler := command.SubmitTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy := "Test User"

	qry := query.GetContractTemplateByIdQuery{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	queryHandler := query.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, template_state.Submitted, contractTemplate.State)
}

func TestSubmit_AllAssigneesApprovedContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Submitted, db)

	ctx := context.Background()
	submittedBy := "Test User"

	assignees := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	ctxTx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = createReviewTasks(t, ctxTx, db, *did, submittedBy, assignees)
	if err != nil {
		t.Fatalf("Failed to create review tasks: %v", err)
	}

	actionFlag := action_flag.Approval

	for _, assignee := range assignees {
		cmd := command.SubmitTemplateContractCommand{
			DID:            *did,
			SubmittedBy:    assignee,
			ActionFlag:     &actionFlag,
			ReviewComments: []string{},
		}
		handler := command.SubmitTemplateContractHandler{
			Ctx: ctx,
			DB:  db,
		}
		err = handler.Handle(cmd)
		if err != nil {
			t.Fatalf("Failed to submit template contract: %v", err)
		}
	}

	retrievedBy := "Test User"

	qry := query.GetContractTemplateByIdQuery{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	queryHandler := query.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, template_state.Reviewed, contractTemplate.State)
}

func TestSubmit_OneAssigneeDeclinesContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Submitted, db)

	ctx := context.Background()
	submittedBy := "Test User"

	assignees := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	ctxTx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = createReviewTasks(t, ctxTx, db, *did, submittedBy, assignees)
	if err != nil {
		t.Fatalf("Failed to create review tasks: %v", err)
	}

	actionFlag := action_flag.Draft

	cmd := command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    assignees[1],
		ActionFlag:     &actionFlag,
		ReviewComments: []string{},
	}
	handler := command.SubmitTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy := "Test User"

	qry := query.GetContractTemplateByIdQuery{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	queryHandler := query.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, template_state.Rejected, contractTemplate.State)
}

func TestSubmit_CompleteWithCreateRejectAndApproveTheContractTemplate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	/**
	Create and Submit the Draft
	*/
	createTestContractTemplate(t, did, template_state.Draft, db)

	ctx := context.Background()
	submittedBy := "Test User"

	assignees := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	cmd := command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    submittedBy,
		ActionFlag:     nil,
		ReviewComments: nil,
		Assignees:      assignees,
	}
	handler := command.SubmitTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy := "Test User"

	qry := query.GetContractTemplateByIdQuery{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	queryHandler := query.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, template_state.Submitted, contractTemplate.State)

	queryReviewTasks := query.GetAllContractTemplateReviewTasksForDID{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	handlerReviewTasks := query.GetAllContractTemplateReviewTasksForDIDHandler{
		Ctx: ctx,
		DB:  db,
	}
	reviewTasks, err := handlerReviewTasks.Handle(queryReviewTasks)
	if err != nil {
		t.Fatalf("Failed to query template review tasks: %v", err)
	}

	for _, reviewTask := range reviewTasks {
		assert.Equal(t, review_task_state.Open, reviewTask.State)

		if !slices.Contains(cmd.Assignees, reviewTask.Assignee) {
			t.Fatalf("Assignee not found in review tasks: %v", reviewTask)
		}
	}

	/**
	First assignee approves the Contract Template
	*/
	actionFlag := action_flag.Approval

	submittedBy = "Test User 1"
	cmd = command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    submittedBy,
		ActionFlag:     &actionFlag,
		ReviewComments: []string{},
	}
	handler = command.SubmitTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy = "Test User"

	qry = query.GetContractTemplateByIdQuery{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	queryHandler = query.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, template_state.Submitted, contractTemplate.State)

	/**
	Second assignee declined the Contract Template
	*/
	actionFlag = action_flag.Draft

	cmd = command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    assignees[1],
		ActionFlag:     &actionFlag,
		ReviewComments: []string{},
	}
	handler = command.SubmitTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy = "Test User"

	qry = query.GetContractTemplateByIdQuery{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	queryHandler = query.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, template_state.Rejected, contractTemplate.State)

	/**
	Contract Template creator submits it again
	*/
	submittedBy = "Test User"
	cmd = command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    submittedBy,
		ActionFlag:     nil,
		ReviewComments: nil,
		Assignees:      assignees,
	}
	handler = command.SubmitTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy = "Test User"

	qry = query.GetContractTemplateByIdQuery{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	queryHandler = query.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, template_state.Submitted, contractTemplate.State)

	queryReviewTasks = query.GetAllContractTemplateReviewTasksForDID{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	handlerReviewTasks = query.GetAllContractTemplateReviewTasksForDIDHandler{
		Ctx: ctx,
		DB:  db,
	}
	reviewTasks, err = handlerReviewTasks.Handle(queryReviewTasks)
	if err != nil {
		t.Fatalf("Failed to query template review tasks: %v", err)
	}

	assert.Equal(t, len(reviewTasks), 6)

	/**
	All assignees approve the Contract Template
	*/
	actionFlag = action_flag.Approval

	for _, assignee := range assignees {
		cmd := command.SubmitTemplateContractCommand{
			DID:            *did,
			SubmittedBy:    assignee,
			ActionFlag:     &actionFlag,
			ReviewComments: []string{},
		}
		handler := command.SubmitTemplateContractHandler{
			Ctx: ctx,
			DB:  db,
		}
		err = handler.Handle(cmd)
		if err != nil {
			t.Fatalf("Failed to submit template contract: %v", err)
		}
	}

	retrievedBy = "Test User"

	qry = query.GetContractTemplateByIdQuery{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	queryHandler = query.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, template_state.Reviewed, contractTemplate.State)
}

func TestSubmit_SubmitContractTemplateInSubmittedStateWithoutActionFlag(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Submitted, db)

	ctx := context.Background()
	submittedBy := "Test User"

	cmd := command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    submittedBy,
		ActionFlag:     nil,
		ReviewComments: []string{},
	}
	handler := command.SubmitTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestSubmit_SubmitContractTemplateInReviewedStateToResubmission(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Reviewed, db)

	ctx := context.Background()
	submittedBy := "Test User"

	cmd := command.SubmitTemplateContractCommand{
		DID:            *did,
		SubmittedBy:    submittedBy,
		ActionFlag:     nil,
		ReviewComments: []string{},
	}
	handler := command.SubmitTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy := "Test User"

	qry := query.GetContractTemplateByIdQuery{
		DID:         *did,
		RetrievedBy: retrievedBy,
	}
	queryHandler := query.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, template_state.Submitted, contractTemplate.State)
}
