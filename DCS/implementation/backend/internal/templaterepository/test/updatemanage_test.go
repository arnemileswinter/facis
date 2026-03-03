package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Draft, creator)

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Submitted, creator)

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Rejected, creator)

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Reviewed, creator)

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Approved, creator)

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Registered, creator)

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Archived, creator)

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	cmd := command.UpdateManageCmd{
		DID:            *did,
		DocumentNumber: 2,
		Version:        2,
		UpdatedAt:      time.Now(),
		UpdatedBy:      "Test User 1",
	}
	handler := command.UpdateManager{
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Draft, creator)
	newState := templatestate.Draft

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Draft, creator)
	newState := templatestate.Submitted

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Draft, creator)
	newState := templatestate.Rejected

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Draft, creator)
	newState := templatestate.Reviewed

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Draft, creator)
	newState := templatestate.Approved

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Draft, creator)
	newState := templatestate.Registered

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Draft, creator)
	newState := templatestate.Archived

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Submitted, creator)
	newState := templatestate.Draft

	reviewers := []string{"Test User 1", "Test User 2", "Test User 3"}
	createReviewTasks(t, ctx, db, repo, *did, reviewtaskstate.Open, creator, reviewers)
	createApprovalTasks(t, ctx, db, repo, *did, approvaltaskstate.Open, creator, "Test User 4")

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
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

	reviewTasksExist, err := repo.RTRepo.TaskExist(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		t.Fatalf("could not check existing review tasks: %v", err)
	}

	approvalTaskExists, err := repo.ATRepo.TaskExists(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Reviewed, creator)
	newState := templatestate.Draft

	reviewers := []string{"Test User 1", "Test User 2", "Test User 3"}
	createReviewTasks(t, ctx, db, repo, *did, reviewtaskstate.Open, creator, reviewers)
	createApprovalTasks(t, ctx, db, repo, *did, approvaltaskstate.Open, creator, "Test User 4")

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
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

	reviewTasksExist, err := repo.RTRepo.TaskExist(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		t.Fatalf("could not check existing review tasks: %v", err)
	}

	approvalTaskExists, err := repo.ATRepo.TaskExists(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Approved, creator)
	newState := templatestate.Draft

	reviewers := []string{"Test User 1", "Test User 2", "Test User 3"}
	createReviewTasks(t, ctx, db, repo, *did, reviewtaskstate.Open, creator, reviewers)
	createApprovalTasks(t, ctx, db, repo, *did, approvaltaskstate.Open, creator, "Test User 4")

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Reviewed, creator)
	newState := templatestate.Submitted

	reviewers := []string{"Test User 1", "Test User 2", "Test User 3"}
	createReviewTasks(t, ctx, db, repo, *did, reviewtaskstate.Approved, creator, reviewers)
	createApprovalTasks(t, ctx, db, repo, *did, approvaltaskstate.Open, creator, "Test User 4")

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
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

	reviewTasksExist, err := repo.RTRepo.ReadAllByID(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Approved, creator)
	newState := templatestate.Submitted

	reviewers := []string{"Test User 1", "Test User 2", "Test User 3"}
	createReviewTasks(t, ctx, db, repo, *did, reviewtaskstate.Approved, creator, reviewers)

	approver := "Test User 4"
	createApprovalTasks(t, ctx, db, repo, *did, approvaltaskstate.Approved, creator, approver)

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Approved, creator)
	newState := templatestate.Reviewed

	reviewers := []string{"Test User 1", "Test User 2", "Test User 3"}
	createReviewTasks(t, ctx, db, repo, *did, reviewtaskstate.Approved, creator, reviewers)

	approver := "Test User 4"
	createApprovalTasks(t, ctx, db, repo, *did, approvaltaskstate.Approved, creator, approver)

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Registered, creator)
	newState := templatestate.Submitted

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Registered, creator)
	newState := templatestate.Submitted

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Registered, creator)
	newState := templatestate.Approved

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Registered, creator)
	newState := templatestate.Archived

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Archived, creator)
	newState := templatestate.Submitted

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Archived, creator)
	newState := templatestate.Submitted

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Archived, creator)
	newState := templatestate.Approved

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
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

	tmpCtx := context.Background()
	ctx, cancel := context.WithTimeout(tmpCtx, base.TransactionTimeout())
	defer cancel()

	repo := NewTestRepo(ctx)

	createContractTemplate(t, db, repo, did, templatestate.Archived, creator)
	newState := templatestate.Registered

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
		Ctx:    ctx,
		DB:     db,
		CTRepo: repo.CTRepo,
		RTRepo: repo.RTRepo,
		ATRepo: repo.ATRepo,
	}
	err = handler.Handle(cmd)

	assert.NotNil(t, err)
}
