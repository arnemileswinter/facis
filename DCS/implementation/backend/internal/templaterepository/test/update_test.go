package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
	"digital-contracting-service/internal/templaterepository/reviewtask"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdate_UpdateContractTemplateDataInDraftState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Draft, creator)

	ctx := context.Background()
	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	retrievedBy := "Test User"

	qry := contracttemplate.GetByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query contract template: %v", err)
	}

	assert.Equal(t, *did, contractTemplate.DID)
	assert.Equal(t, name, *contractTemplate.Name)
	assert.Equal(t, description, *contractTemplate.Description)
	//assert.Equal(t, jsonTemplateData, contractTemplate.TemplateData)
}

func TestUpdate_UpdateNonExistingContractTemplate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	ctx := context.Background()

	cmd := command.UpdateCommand{
		DID:            *did,
		DocumentNumber: 2,
		Version:        2,
		UpdatedAt:      time.Now(),
		UpdatedBy:      "Test User 1",
	}
	handler := command.UpdateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateDataInDraftStateWithInvalidUser(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Draft, creator)

	reviewers := []string{"Test User 2"}

	ctx := context.Background()
	createReviewTasks(t, ctx, db, *did, reviewtaskstate.Open, creator, reviewers)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      "Test User 1",
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateDataInSubmittedStateAsCreator(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Submitted, creator)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateDataInSubmittedStateAsReviewer(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Submitted, creator)

	reviewers := []string{"Test User 2"}

	ctx := context.Background()
	createReviewTasks(t, ctx, db, *did, reviewtaskstate.Open, creator, reviewers)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      reviewers[0],
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	retrievedBy := "Test User"

	qry := contracttemplate.GetByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query contract template: %v", err)
	}

	assert.Equal(t, *did, contractTemplate.DID)
	assert.Equal(t, name, *contractTemplate.Name)
	assert.Equal(t, description, *contractTemplate.Description)
	//assert.Equal(t, jsonTemplateData, contractTemplate.TemplateData)
}

func TestUpdate_UpdateContractTemplateDataInSubmittedStateWithInvalidUser(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Submitted, creator)

	reviewers := []string{"Test User 2"}

	ctx := context.Background()
	createReviewTasks(t, ctx, db, *did, reviewtaskstate.Open, creator, reviewers)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      "Test User 1",
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateDataInDraftApprovedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Approved, creator)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateDataInDraftPublishedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Registered, creator)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateDataInDraftArchivedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Archived, creator)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateDataInDraftApprovedStateWithInvalidUser(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Approved, creator)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      "Test User 1",
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateHandler{
		Ctx: context.Background(),
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateAfterUpdate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Draft, creator)

	ctx := context.Background()
	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now().Add(-5 * time.Second),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdate_UpdateContractTemplateAndReopenTasks(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Submitted, creator)

	ctx := context.Background()

	ctxTx, cancel := context.WithTimeout(ctx, base.TransactionTimeout())
	defer cancel()

	reviewers := []string{
		"Test User 1",
		"Test User 2",
		"Test User 3",
	}

	createReviewTasks(t, ctxTx, db, *did, reviewtaskstate.Approved, creator, reviewers)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateCommand{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      reviewers[1],
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to update template: %v", err)
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		t.Fatal("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	exists, err := reviewtask.ExistTasksInStates(ctx, tx, *did, 1, 1, reviewtaskstate.Approved, reviewtaskstate.Verified, reviewtaskstate.Rejected)
	if err != nil {
		t.Fatalf("Failed to check existence of review tasks: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal("could not commit transaction: %w", err)
	}

	assert.False(t, exists)
}
