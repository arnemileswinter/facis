package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/template_repository"
	"digital-contracting-service/internal/template_repository/command"
	"digital-contracting-service/internal/template_repository/datatype/action_flag"
	"digital-contracting-service/internal/template_repository/datatype/approval_task_state"
	"digital-contracting-service/internal/template_repository/datatype/review_task_state"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"digital-contracting-service/internal/template_repository/query"
	"slices"
	"testing"

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
	approver := "Test User 5"
	cmd := command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		SubmittedBy:    submittedBy,
		ActionFlag:     nil,
		Comments:       nil,
		Reviewer: []string{
			"Test User 2",
			"Test User 3",
			"Test User 4",
		},
		Approver: &approver,
	}
	handler := command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

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

	assert.Equal(t, template_state.Submitted, contractTemplate.State)

	queryReviewTasks := query.GetAllContractTemplateReviewTasksForDID{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
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

		if !slices.Contains(cmd.Reviewer, reviewTask.Reviewer) {
			t.Fatalf("Reviewer not found in review tasks: %v", reviewTask)
		}
	}
}

func createReviewTasks(t *testing.T, ctx context.Context, db *sqlx.DB, did string, submittedBy string, reviewers []string) error {
	tx, err := db.BeginTxx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	for _, reviewer := range reviewers {
		reviewTask := template_repository.ReviewTaskData{
			DID:            did,
			DocumentNumber: 1,
			Version:        1,
			Reviewer:       reviewer,
			State:          review_task_state.Open,
			CreatedBy:      submittedBy,
		}
		_, err = template_repository.CreateReviewTasks(ctx, tx, reviewTask)
		if err != nil {
			t.Fatalf("Failed to create review task: %v", err)
		}
	}
	return tx.Commit()
}

func TestSubmit_OneReviewerApprovedContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Submitted, db)

	ctx := context.Background()
	submittedBy := "Test User"

	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	ctxTx, cancel := context.WithTimeout(ctx, base.GetTransactionTimeout())
	defer cancel()
	err = createReviewTasks(t, ctxTx, db, *did, submittedBy, reviewers)
	if err != nil {
		t.Fatalf("Failed to create review tasks: %v", err)
	}

	tx, err := db.BeginTxx(ctxTx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	actionFlag := action_flag.Approval

	submittedBy = "Test User 1"
	cmd := command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		SubmittedBy:    submittedBy,
		ActionFlag:     &actionFlag,
		Comments:       []string{},
	}
	handler := command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

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

	assert.Equal(t, template_state.Submitted, contractTemplate.State)
}

func TestSubmit_AllReviewersApprovedContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Submitted, db)

	ctx := context.Background()
	submittedBy := "Test User"

	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	ctxTx, cancel := context.WithTimeout(ctx, base.GetTransactionTimeout())
	defer cancel()
	err = createReviewTasks(t, ctxTx, db, *did, submittedBy, reviewers)
	if err != nil {
		t.Fatalf("Failed to create review tasks: %v", err)
	}

	tx, err := db.BeginTxx(ctxTx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	actionFlag := action_flag.Approval

	for _, reviewer := range reviewers {
		cmd := command.SubmitContractTemplateCommand{
			DID:            *did,
			DocumentNumber: 1,
			Version:        1,
			SubmittedBy:    reviewer,
			ActionFlag:     &actionFlag,
			Comments:       []string{},
		}
		handler := command.SubmitContractTemplateHandler{
			Ctx: ctx,
			DB:  db,
		}
		err = handler.Handle(cmd)
		if err != nil {
			t.Fatalf("Failed to submit template contract: %v", err)
		}
	}

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

	assert.Equal(t, template_state.Reviewed, contractTemplate.State)
}

func TestSubmit_OneReviewerDeclinesContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, template_state.Submitted, db)

	ctx := context.Background()

	submittedBy := "Test User"

	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	ctxTx, cancel := context.WithTimeout(ctx, base.GetTransactionTimeout())
	defer cancel()
	err = createReviewTasks(t, ctxTx, db, *did, submittedBy, reviewers)
	if err != nil {
		t.Fatalf("Failed to create review tasks: %v", err)
	}

	tx, err := db.BeginTxx(ctxTx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	actionFlag := action_flag.Draft

	cmd := command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		SubmittedBy:    reviewers[1],
		ActionFlag:     &actionFlag,
		Comments:       []string{},
	}
	handler := command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

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

	assert.Equal(t, template_state.Rejected, contractTemplate.State)
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

	cmd := command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		SubmittedBy:    submittedBy,
		ActionFlag:     nil,
		Comments:       []string{},
	}
	handler := command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestSubmit_SubmitContractTemplateWithResubmission(t *testing.T) {

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
	approver := "Test User 4"
	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	cmd := command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		SubmittedBy:    submittedBy,
		ActionFlag:     nil,
		Comments:       nil,
		Reviewer:       reviewers,
		Approver:       &approver,
	}
	handler := command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

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

	assert.Equal(t, template_state.Submitted, contractTemplate.State)

	queryReviewTasks := query.GetAllContractTemplateReviewTasksForDID{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	handlerReviewTasks := query.GetAllContractTemplateReviewTasksForDIDHandler{
		Ctx: ctx,
		DB:  db,
	}
	reviewTasks, err := handlerReviewTasks.Handle(queryReviewTasks)
	if err != nil {
		t.Fatalf("Failed to query template review tasks: %v", err)
	}

	assert.Equal(t, len(reviewTasks), 3)

	for _, reviewTask := range reviewTasks {
		assert.Equal(t, review_task_state.Open, reviewTask.State)

		if !slices.Contains(cmd.Reviewer, reviewTask.Reviewer) {
			t.Fatalf("Reviewer not found in review tasks: %v", reviewTask)
		}
	}

	queryApprovalTasks := query.GetAllContractTemplateApprovalTasksForDID{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	handlerApprovalTasks := query.GetAllContractTemplateApprovalTasksForDIDHandler{
		Ctx: ctx,
		DB:  db,
	}
	approvalTasks, err := handlerApprovalTasks.Handle(queryApprovalTasks)
	if err != nil {
		t.Fatalf("Failed to query template review tasks: %v", err)
	}

	assert.Equal(t, len(approvalTasks), 1)
	assert.Equal(t, approval_task_state.Open, approvalTasks[0].State)

	/**
	First reviewer approves the Contract Template
	*/
	actionFlag := action_flag.Approval

	submittedBy = "Test User 1"
	cmd = command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		SubmittedBy:    submittedBy,
		ActionFlag:     &actionFlag,
		Comments:       []string{},
	}
	handler = command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy = "Test User"

	qry = query.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
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
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	handlerReviewTasks = query.GetAllContractTemplateReviewTasksForDIDHandler{
		Ctx: ctx,
		DB:  db,
	}
	reviewTasks, err = handlerReviewTasks.Handle(queryReviewTasks)
	if err != nil {
		t.Fatalf("Failed to query template review tasks: %v", err)
	}

	assert.Equal(t, len(reviewTasks), 3)

	queryApprovalTasks = query.GetAllContractTemplateApprovalTasksForDID{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	handlerApprovalTasks = query.GetAllContractTemplateApprovalTasksForDIDHandler{
		Ctx: ctx,
		DB:  db,
	}
	approvalTasks, err = handlerApprovalTasks.Handle(queryApprovalTasks)
	if err != nil {
		t.Fatalf("Failed to query template review tasks: %v", err)
	}

	assert.Equal(t, len(approvalTasks), 1)

	/**
	Second reviewer declined the Contract Template
	*/
	actionFlag = action_flag.Draft

	cmd = command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		SubmittedBy:    reviewers[1],
		ActionFlag:     &actionFlag,
		Comments:       []string{},
	}
	handler = command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy = "Test User"

	qry = query.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
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
	approver = "Test User 4"
	reviewers = []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	submittedBy = "Test User"
	cmd = command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		SubmittedBy:    submittedBy,
		ActionFlag:     nil,
		Comments:       nil,
		Approver:       &approver,
		Reviewer:       reviewers,
	}
	handler = command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy = "Test User"

	qry = query.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
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
	All reviewer approve the Contract Template
	*/
	actionFlag = action_flag.Approval

	for _, reviewer := range reviewers {
		cmd := command.SubmitContractTemplateCommand{
			DID:            *did,
			DocumentNumber: 1,
			Version:        1,
			SubmittedBy:    reviewer,
			ActionFlag:     &actionFlag,
			Comments:       []string{},
		}
		handler := command.SubmitContractTemplateHandler{
			Ctx: ctx,
			DB:  db,
		}
		err = handler.Handle(cmd)
		if err != nil {
			t.Fatalf("Failed to submit template contract: %v", err)
		}
	}

	retrievedBy = "Test User"

	qry = query.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
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

	/**
	Approver resubmits reviewed Contract Template
	*/
	cmd = command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		SubmittedBy:    approver,
		ActionFlag:     nil,
		Comments:       []string{"Test Comment"},
		Reviewer:       nil,
	}
	handler = command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy = "Test User"

	qry = query.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
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
	All reviewer approve the Contract Template
	*/
	actionFlag = action_flag.Approval

	for _, reviewer := range reviewers {
		cmd := command.SubmitContractTemplateCommand{
			DID:            *did,
			DocumentNumber: 1,
			Version:        1,
			SubmittedBy:    reviewer,
			ActionFlag:     &actionFlag,
			Comments:       []string{},
		}
		handler := command.SubmitContractTemplateHandler{
			Ctx: ctx,
			DB:  db,
		}
		err = handler.Handle(cmd)
		if err != nil {
			t.Fatalf("Failed to submit template contract: %v", err)
		}
	}

	retrievedBy = "Test User"

	qry = query.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
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

	/**
	Approver resubmits reviewed Contract Template
	*/
	cmd = command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		SubmittedBy:    approver,
		ActionFlag:     nil,
		Comments:       []string{"Test Comment"},
		Reviewer:       nil,
	}
	handler = command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy = "Test User"

	qry = query.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
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
}
