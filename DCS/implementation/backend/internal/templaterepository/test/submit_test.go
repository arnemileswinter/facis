package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/actionflag"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/query"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
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

	createTestContractTemplate(t, did, templatestate.Draft, db)

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

	qry := contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

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
		assert.Equal(t, reviewtaskstate.Open, reviewTask.State)

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
		reviewTask := templaterepository.ReviewTaskData{
			DID:            did,
			DocumentNumber: 1,
			Version:        1,
			Reviewer:       reviewer,
			State:          reviewtaskstate.Open,
			CreatedBy:      submittedBy,
		}
		_, err = templaterepository.CreateReviewTasks(ctx, tx, reviewTask)
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

	createTestContractTemplate(t, did, templatestate.Submitted, db)

	ctx := context.Background()
	submittedBy := "Test User"

	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
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

	actionFlag := actionflag.Approval

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

	qry := contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)
}

func TestSubmit_AllReviewersApprovedContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, templatestate.Submitted, db)

	ctx := context.Background()
	submittedBy := "Test User"

	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
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

	actionFlag := actionflag.Approval

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

	qry := contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Reviewed, contractTemplate.State)
}

func TestSubmit_OneReviewerDeclinesContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, templatestate.Submitted, db)

	ctx := context.Background()

	submittedBy := "Test User"

	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
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

	actionFlag := actionflag.Draft

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

	qry := contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Rejected, contractTemplate.State)
}

func TestSubmit_SubmitContractTemplateInSubmittedStateWithoutActionFlag(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, templatestate.Submitted, db)

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
	createTestContractTemplate(t, did, templatestate.Draft, db)

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

	qry := contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

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
		assert.Equal(t, reviewtaskstate.Open, reviewTask.State)

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
	assert.Equal(t, approvaltaskstate.Open, approvalTasks[0].State)

	/**
	First reviewer approves the Contract Template
	*/
	actionFlag := actionflag.Approval

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

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
	actionFlag = actionflag.Draft

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Rejected, contractTemplate.State)

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

	/**
	All reviewer approve the Contract Template
	*/
	actionFlag = actionflag.Approval

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Reviewed, contractTemplate.State)

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

	/**
	All reviewer approve the Contract Template
	*/
	actionFlag = actionflag.Approval

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Reviewed, contractTemplate.State)

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)
}

func TestSubmit_SubmitContractTemplateWithApproving(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	/**
	Create and Submit the Draft
	*/
	createTestContractTemplate(t, did, templatestate.Draft, db)

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

	qry := contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

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
		assert.Equal(t, reviewtaskstate.Open, reviewTask.State)

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
	assert.Equal(t, approvaltaskstate.Open, approvalTasks[0].State)

	/**
	First reviewer approves the Contract Template
	*/
	actionFlag := actionflag.Approval

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

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
	actionFlag = actionflag.Draft

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Rejected, contractTemplate.State)

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

	/**
	All reviewer approve the Contract Template
	*/
	actionFlag = actionflag.Approval

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Reviewed, contractTemplate.State)

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

	/**
	All reviewer approve the Contract Template
	*/
	actionFlag = actionflag.Approval

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Reviewed, contractTemplate.State)

	/**
	Approver approves reviewed Contract Template
	*/
	approveCmd := command.ApproveTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		ApprovedBy:     approver,
		DecisionNotes:  []string{"Test"},
	}
	approveHandler := command.ApproveTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = approveHandler.Handle(approveCmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy = "Test User"

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Approved, contractTemplate.State)
}

func TestSubmit_SubmitContractTemplateWithRejecting(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	/**
	Create and Submit the Draft
	*/
	createTestContractTemplate(t, did, templatestate.Draft, db)

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

	qry := contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

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
		assert.Equal(t, reviewtaskstate.Open, reviewTask.State)

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
	assert.Equal(t, approvaltaskstate.Open, approvalTasks[0].State)

	/**
	First reviewer approves the Contract Template
	*/
	actionFlag := actionflag.Approval

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

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
	actionFlag = actionflag.Draft

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Rejected, contractTemplate.State)

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

	/**
	All reviewer approve the Contract Template
	*/
	actionFlag = actionflag.Approval

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Reviewed, contractTemplate.State)

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

	/**
	All reviewer approve the Contract Template
	*/
	actionFlag = actionflag.Approval

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

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Reviewed, contractTemplate.State)

	/**
	Approver rejects reviewed Contract Template
	*/
	rejectCmd := command.RejectTemplateContractCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RejectedBy:     approver,
		Reason:         "Test",
	}
	rejectHandler := command.RejectTemplateContractHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = rejectHandler.Handle(rejectCmd)
	if err != nil {
		t.Fatalf("Failed to submit template contract: %v", err)
	}

	retrievedBy = "Test User"

	qry = contracttemplate.GetContractTemplatesByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler = contracttemplate.GetContractTemplateByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err = queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, templatestate.Draft, contractTemplate.State)
}
