package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/datatype/templatetype"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch_SearchContractTemplatesWithoutSearchValue(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	creator := "Test User"

	templateData := map[string]interface{}{}

	did, _ := base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 1, 3, "Test1", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "Test1", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "Test1", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "Test1", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 3, 2, "Test1", "Test1", templateData)

	ctx := context.Background()

	qry := contracttemplate.GetAllMetaDataByFilterQuery{
		RetrievedBy: creator,
	}
	queryHandler := contracttemplate.GetAllMetaDataByFilterHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, 5, len(contractTemplate))
}

func TestSearch_SearchContractTemplatesByDID(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	name := "Test Contract Template"
	description := "Test Description"

	templateData := map[string]interface{}{}
	jsonMetaData, err := datatype.NewJSON(templateData)
	if err != nil {
		t.Fatalf("Failed to create JSON data: %v", err)
	}

	creator := "Test User"

	ctx := context.Background()

	cmd := command.CreateCommand{
		DID:          *did,
		CreatedBy:    creator,
		TemplateType: templatetype.FrameContract,
		Name:         &name,
		Description:  &description,
		TemplateData: &jsonMetaData,
	}
	createHandler := command.CreateHandler{
		Ctx: ctx,
		DB:  db,
	}
	err = createHandler.Handle(cmd)
	if err != nil {
		t.Fatalf("Failed to create template contract: %v", err)
	}

	qry := contracttemplate.GetAllMetaDataByFilterQuery{
		RetrievedBy: creator,
		DID:         did,
	}
	queryHandler := contracttemplate.GetAllMetaDataByFilterHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, 1, len(contractTemplate))
}

func TestSearch_SearchContractTemplatesByDocumentNumberAndVersion(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	creator := "Test User"

	templateData := map[string]interface{}{}

	did, _ := base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 1, 3, "Test1", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "Test1", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "Test1", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "Test1", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 3, 2, "Test1", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 3, 2, "Test1", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 3, 2, "Test1", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 3, 1, "Test1", "Test1", templateData)

	ctx := context.Background()

	documentNumber := 3
	version := 2
	qry := contracttemplate.GetAllMetaDataByFilterQuery{
		RetrievedBy:    creator,
		DocumentNumber: &documentNumber,
		Version:        &version,
	}
	queryHandler := contracttemplate.GetAllMetaDataByFilterHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, 3, len(contractTemplate))
}

func TestSearch_SearchContractTemplatesByName(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	creator := "Test User"

	templateData := map[string]interface{}{}

	did, _ := base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 1, 3, "-- test 1 --", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "-- test 1.2 --", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "-- test 1.3 --", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "-- test 2 --", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 3, 2, "-- test 2.2 --", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 3, 2, "-- test 2.3 --", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 3, 2, "-- test 3 --", "Test1", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 3, 1, "-- test 3.2 --", "Test1", templateData)

	ctx := context.Background()

	searchName := "Test 2." // The search is case-insensitive
	qry := contracttemplate.GetAllMetaDataByFilterQuery{
		RetrievedBy: creator,
		Name:        &searchName,
	}
	queryHandler := contracttemplate.GetAllMetaDataByFilterHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, 2, len(contractTemplate))
}

func TestSearch_SearchContractTemplatesByDescript(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	creator := "Test User"

	templateData := map[string]interface{}{}

	did, _ := base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 1, 3, "-- test 1 --", "a long test1 description", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "-- test 1.2 --", "a long test2 description", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "-- test 1.3 --", "a long test2.2 description", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "-- test 2 --", "a long test2.3 description", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 3, 2, "-- test 2.2 --", "a long test3 description", templateData)

	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 3, 2, "-- test 2.3 --", "a long test4 description", templateData)

	ctx := context.Background()

	searchDescription := "Test2." // The search is case-insensitive
	qry := contracttemplate.GetAllMetaDataByFilterQuery{
		RetrievedBy: creator,
		Description: &searchDescription,
	}
	queryHandler := contracttemplate.GetAllMetaDataByFilterHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, 2, len(contractTemplate))
}

func TestSearch_SearchContractTemplatesByTemplateData(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	creator := "Test User"

	templateData := map[string]interface{}{
		"name":        "-- test1 --",
		"description": "a long test1 description",
	}
	did, _ := base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 1, 3, "-- test1 --", "a long test1 description", templateData)

	templateData = map[string]interface{}{
		"name":        "-- test1.2 --",
		"description": "a long test2 description",
	}
	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "-- test1.2 --", "a long test2 description", templateData)

	templateData = map[string]interface{}{
		"name":        "-- test1.3 --",
		"description": "a long test2.2 description",
	}
	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "-- test1.3 --", "a long test2.2 description", templateData)

	templateData = map[string]interface{}{
		"name":        "-- test2 --",
		"description": "a long test2.3 description",
	}
	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 2, 3, "-- test2 --", "a long test2.3 description", templateData)

	templateData = map[string]interface{}{
		"name":        "-- test2.2 --",
		"description": "a long test3 description",
	}
	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 3, 2, "-- test2.2 --", "a long test3 description", templateData)

	templateData = map[string]interface{}{
		"name":        "-- test2.3 --",
		"description": "a long test4 description",
	}
	did, _ = base.GetDID()
	createTestContractTemplateWithData(t, db, did, templatestate.Reviewed, creator, 3, 2, "-- test2.3 --", "a long test4 description", templateData)

	ctx := context.Background()

	filter := "Test2.2" // The search is case-insensitive
	qry := contracttemplate.GetAllMetaDataByFilterQuery{
		RetrievedBy: creator,
		Filter:      &filter,
	}
	queryHandler := contracttemplate.GetAllMetaDataByFilterHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, 2, len(contractTemplate))
}
