package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/approvaltask"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
	"digital-contracting-service/internal/templaterepository/reviewtask"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateManage_UpdateContractTemplateDataInDraftState(t *testing.T) {

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	retrievedBy := "Test User"

	qry := contracttemplate.GetByIDQry{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetByIDHandler{
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

func TestUpdateManage_UpdateContractTemplateDataInSubmitState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Submitted, creator)

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	retrievedBy := "Test User"

	qry := contracttemplate.GetByIDQry{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetByIDHandler{
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

func TestUpdateManage_UpdateContractTemplateDataInRejectedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Rejected, creator)

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	retrievedBy := "Test User"

	qry := contracttemplate.GetByIDQry{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetByIDHandler{
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

func TestUpdateManage_UpdateContractTemplateDataInReviewedState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Reviewed, creator)

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	retrievedBy := "Test User"

	qry := contracttemplate.GetByIDQry{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetByIDHandler{
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

func TestUpdateManage_UpdateContractTemplateDataInApproveState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Approved, creator)

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_UpdateContractTemplateDataInRegisteredState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Registered, creator)

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_UpdateContractTemplateDataInArchiveState(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Archived, creator)

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_UpdateNonExistingContractTemplate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	ctx := context.Background()

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 2,
		Version:        2,
		UpdatedAt:      time.Now(),
		UpdatedBy:      "Test User 1",
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromDraftToDraft(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Draft, creator)
	newState := templatestate.Draft

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	retrievedBy := "Test User"

	qry := contracttemplate.GetByIDQry{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetByIDHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query contract template: %v", err)
	}

	assert.Equal(t, *did, contractTemplate.DID)
	assert.Equal(t, newState, contractTemplate.State)
}

func TestUpdateManage_SetContractTemplateStateFromDraftToSubmitted(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Draft, creator)
	newState := templatestate.Submitted

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromDraftToRejected(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Draft, creator)
	newState := templatestate.Rejected

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromDraftToReviewed(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Draft, creator)
	newState := templatestate.Reviewed

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromDraftToApproved(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Draft, creator)
	newState := templatestate.Approved

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromDraftToRegistered(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Draft, creator)
	newState := templatestate.Registered

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromDraftToArchive(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Draft, creator)
	newState := templatestate.Archived

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	retrievedBy := "Test User"

	qry := contracttemplate.GetByIDQry{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetByIDHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query contract template: %v", err)
	}

	assert.Equal(t, *did, contractTemplate.DID)
	assert.Equal(t, newState, contractTemplate.State)
}

func TestUpdateManage_SetContractTemplateStateFromSubmittedToDraft(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Submitted, creator)
	newState := templatestate.Draft

	ctx := context.Background()

	reviewers := []string{"Test User 1", "Test User 2", "Test User 3"}
	createReviewTasks(t, ctx, db, *did, reviewtaskstate.Open, creator, reviewers)
	createApprovalTasks(t, ctx, db, *did, approvaltaskstate.Open, creator, "Test User 4")

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	retrievedBy := "Test User"

	qry := contracttemplate.GetByIDQry{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetByIDHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query contract template: %v", err)
	}

	tx, err := db.BeginTxx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	reviewTasksExist, err := reviewtask.TaskExist(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		t.Fatalf("could not check existing review tasks: %v", err)
	}

	approvalTaskExists, err := approvaltask.TaskExists(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		t.Fatalf("could not check existing approval tasks: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	assert.Equal(t, *did, contractTemplate.DID)
	assert.Equal(t, newState, contractTemplate.State)
	assert.False(t, reviewTasksExist)
	assert.False(t, approvalTaskExists)
}

func TestUpdateManage_SetContractTemplateStateFromReviewedToDraft(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Reviewed, creator)
	newState := templatestate.Draft

	ctx := context.Background()

	reviewers := []string{"Test User 1", "Test User 2", "Test User 3"}
	createReviewTasks(t, ctx, db, *did, reviewtaskstate.Open, creator, reviewers)
	createApprovalTasks(t, ctx, db, *did, approvaltaskstate.Open, creator, "Test User 4")

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	retrievedBy := "Test User"

	qry := contracttemplate.GetByIDQry{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetByIDHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query contract template: %v", err)
	}

	tx, err := db.BeginTxx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	reviewTasksExist, err := reviewtask.TaskExist(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		t.Fatalf("could not check existing review tasks: %v", err)
	}

	approvalTaskExists, err := approvaltask.TaskExists(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		t.Fatalf("could not check existing approval tasks: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	assert.Equal(t, *did, contractTemplate.DID)
	assert.Equal(t, newState, contractTemplate.State)
	assert.False(t, reviewTasksExist)
	assert.False(t, approvalTaskExists)
}

func TestUpdateManage_SetContractTemplateStateFromApprovedToDraft(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Approved, creator)
	newState := templatestate.Draft

	ctx := context.Background()

	reviewers := []string{"Test User 1", "Test User 2", "Test User 3"}
	createReviewTasks(t, ctx, db, *did, reviewtaskstate.Open, creator, reviewers)
	createApprovalTasks(t, ctx, db, *did, approvaltaskstate.Open, creator, "Test User 4")

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromReviewedToSubmitted(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Reviewed, creator)
	newState := templatestate.Submitted

	ctx := context.Background()

	reviewers := []string{"Test User 1", "Test User 2", "Test User 3"}
	createReviewTasks(t, ctx, db, *did, reviewtaskstate.Approved, creator, reviewers)
	createApprovalTasks(t, ctx, db, *did, approvaltaskstate.Open, creator, "Test User 4")

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to submit contract template: %v", err)
	}

	retrievedBy := "Test User"

	qry := contracttemplate.GetByIDQry{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    retrievedBy,
	}
	queryHandler := contracttemplate.GetByIDHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query contract template: %v", err)
	}

	tx, err := db.BeginTxx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	reviewTasksExist, err := reviewtask.ReadAllByID(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		t.Fatalf("could not check existing review tasks: %v", err)
	}

	tasksAreOpen := true
	for _, reviewTask := range reviewTasksExist {
		if reviewTask.State != reviewtaskstate.Open {
			tasksAreOpen = false
		}
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	assert.Equal(t, *did, contractTemplate.DID)
	assert.Equal(t, newState, contractTemplate.State)
	assert.True(t, tasksAreOpen)
}

func TestUpdateManage_SetContractTemplateStateFromApprovedToSubmitted(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Approved, creator)
	newState := templatestate.Submitted

	ctx := context.Background()

	reviewers := []string{"Test User 1", "Test User 2", "Test User 3"}
	createReviewTasks(t, ctx, db, *did, reviewtaskstate.Approved, creator, reviewers)

	approver := "Test User 4"
	createApprovalTasks(t, ctx, db, *did, approvaltaskstate.Approved, creator, approver)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromApprovedToReviewed(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Approved, creator)
	newState := templatestate.Reviewed

	ctx := context.Background()

	reviewers := []string{"Test User 1", "Test User 2", "Test User 3"}
	createReviewTasks(t, ctx, db, *did, reviewtaskstate.Approved, creator, reviewers)

	approver := "Test User 4"
	createApprovalTasks(t, ctx, db, *did, approvaltaskstate.Approved, creator, approver)

	templateData := map[string]interface{}{
		"test": "update",
	}
	jsonTemplateData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON template data: %v", err)
	}

	name := "Updated Contract Template"
	description := "Updated Description"

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromRegisteredToDraft(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Registered, creator)
	newState := templatestate.Submitted

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromRegisteredToSubmitted(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Registered, creator)
	newState := templatestate.Submitted

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromRegisteredToApproved(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Registered, creator)
	newState := templatestate.Approved

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromRegisteredToArchived(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Registered, creator)
	newState := templatestate.Archived

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromArchivedToDraft(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Archived, creator)
	newState := templatestate.Submitted

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromArchivedToSubmitted(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Archived, creator)
	newState := templatestate.Submitted

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromArchivedToApproved(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Archived, creator)
	newState := templatestate.Approved

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}

func TestUpdateManage_SetContractTemplateStateFromArchivedToRegistered(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Archived, creator)
	newState := templatestate.Registered

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

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		State:          &newState,
		UpdatedBy:      creator,
		UpdatedAt:      time.Now(),
		Name:           &name,
		Description:    &description,
		TemplateData:   &jsonTemplateData,
	}
	handler := command.UpdateManager{
		Ctx: ctx,
		DB:  db,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}
