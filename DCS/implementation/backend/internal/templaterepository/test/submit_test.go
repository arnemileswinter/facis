package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/actionflag"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/query"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubmit_SubmitContractTemplateInDraftState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Draft, creator)

	ctx := context.Background()

	approver := "Test User 5"
	cmd := command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		SubmittedBy:    creator,
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

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

	queryReviewTasks := query.GetAllContractTemplateReviewTasksForDID{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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

func TestSubmit_SubmitContractTemplateInDraftStateWithInvalidUser(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Draft, creator)

	ctx := context.Background()

	approver := "Test User 5"
	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	cmd := command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		SubmittedBy:    "Test User 6",
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

	assert.NotNil(t, err)
}

func TestSubmit_OneReviewerApprovedContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Submitted, creator)

	ctx := context.Background()

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
	defer cancel()

	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	createReviewTasks(t, ctxTx, db, *did, reviewtaskstate.Open, creator, reviewers)

	createApprovalTasks(t, ctxTx, db, *did, approvaltaskstate.Open, creator, "Test User 4")

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

	cmd := command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		SubmittedBy:    reviewers[0],
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

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)
}

func TestSubmit_ApproveContractTemplateInSubmittedStateWithInvalidUser(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Submitted, creator)

	ctx := context.Background()

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
	defer cancel()

	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	createReviewTasks(t, ctxTx, db, *did, reviewtaskstate.Open, creator, reviewers)

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

	cmd := command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		SubmittedBy:    "Test User 4",
		ActionFlag:     &actionFlag,
		Comments:       []string{},
	}
	handler := command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestSubmit_RejectContractTemplateInSubmittedStateWithInvalidUser(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Submitted, creator)

	ctx := context.Background()

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
	defer cancel()

	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	createReviewTasks(t, ctxTx, db, *did, reviewtaskstate.Open, creator, reviewers)

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
		UpdatedAt:      time.Now(),
		SubmittedBy:    "Test User 4",
		ActionFlag:     &actionFlag,
		Comments:       []string{},
	}
	handler := command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestSubmit_AllReviewersApprovedContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Submitted, creator)

	ctx := context.Background()

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
	defer cancel()

	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	createReviewTasks(t, ctxTx, db, *did, reviewtaskstate.Open, creator, reviewers)

	createApprovalTasks(t, ctxTx, db, *did, approvaltaskstate.Open, creator, "Test User 4")

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
			UpdatedAt:      time.Now(),
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

	assert.Equal(t, templatestate.Reviewed, contractTemplate.State)
}

func TestSubmit_OneReviewerDeclinesContractTemplateInSubmittedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Submitted, creator)

	ctx := context.Background()

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
	defer cancel()

	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	createReviewTasks(t, ctxTx, db, *did, reviewtaskstate.Open, creator, reviewers)

	createApprovalTasks(t, ctxTx, db, *did, approvaltaskstate.Open, creator, "Test User 4")

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
		UpdatedAt:      time.Now(),
		SubmittedBy:    reviewers[0],
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

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Submitted, creator)

	ctx := context.Background()

	cmd := command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		SubmittedBy:    creator,
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

func TestSubmit_SubmitContractTemplateInReviewedStateWithInvalidUser(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	ctx := context.Background()

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
	defer cancel()

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Reviewed, creator)

	approver := "Test User 1"

	createApprovalTasks(t, ctxTx, db, *did, approvaltaskstate.Open, creator, approver)

	cmd := command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		SubmittedBy:    "Test User 2",
	}
	handler := command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.Error(t, err)
}

func TestSubmit_SubmitContractTemplateInSubmittedStateWithApproverUser(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	ctx := context.Background()

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
	defer cancel()

	creator := "Test User"

	createTestContractTemplate(t, db, did, templatestate.Submitted, creator)

	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	createReviewTasks(t, ctxTx, db, *did, reviewtaskstate.Open, creator, reviewers)

	approver := "Test User 4"

	createApprovalTasks(t, ctxTx, db, *did, approvaltaskstate.Open, creator, approver)

	aFlag := actionflag.Approval
	cmd := command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		SubmittedBy:    approver,
		ActionFlag:     &aFlag,
	}
	handler := command.SubmitContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.Error(t, err)
}

func TestSubmit_SubmitContractTemplateWithResubmission(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	creator := "Test User"

	/**
	Create and Submit the Draft
	*/
	createTestContractTemplate(t, db, did, templatestate.Draft, creator)

	ctx := context.Background()

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
		UpdatedAt:      time.Now(),
		SubmittedBy:    creator,
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

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

	queryReviewTasks := query.GetAllContractTemplateReviewTasksForDID{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
		RetrievedBy:    creator,
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

	cmd = command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
		RetrievedBy:    creator,
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
		RetrievedBy:    creator,
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
		UpdatedAt:      time.Now(),
		SubmittedBy:    reviewers[2],
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
	cmd = command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		SubmittedBy:    creator,
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
			UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
		UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
			UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
		UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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

	creator := "Test User"

	/**
	Create and Submit the Draft
	*/
	createTestContractTemplate(t, db, did, templatestate.Draft, creator)

	ctx := context.Background()

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
		UpdatedAt:      time.Now(),
		SubmittedBy:    creator,
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

	assert.Equal(t, templatestate.Submitted, contractTemplate.State)

	queryReviewTasks := query.GetAllContractTemplateReviewTasksForDID{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
		RetrievedBy:    creator,
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

	cmd = command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
		RetrievedBy:    creator,
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
		RetrievedBy:    creator,
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
		UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
	cmd = command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now(),
		SubmittedBy:    creator,
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
			UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
		UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
			UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
		UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
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
	createTestContractTemplate(t, db, did, templatestate.Draft, "Test User")

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
		UpdatedAt:      time.Now(),
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
		UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
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
		UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
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
		UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
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
			UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
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
		UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
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
			UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
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
		UpdatedAt:      time.Now(),
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

	qry = contracttemplate.GetContractTemplateByIdQuery{
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

func TestSubmit_SubmitContractTemplateAfterUpdate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, db, did, templatestate.Draft, "Test User")

	ctx := context.Background()
	submittedBy := "Test User"
	approver := "Test User 5"
	cmd := command.SubmitContractTemplateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedAt:      time.Now().Add(-5 * time.Minute),
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

	assert.NotNil(t, err)
}
