package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetrieve_RetrieveContractTemplateById(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to get new DID: %v", err)
	}

	creator := "Test User"

	createContractTemplate(t, db, did, templatestate.Draft, creator)

	ctx := context.Background()

	qry := contracttemplate.GetByIdQuery{
		DID:            *did,
		DocumentNumber: 1,
		Version:        1,
		RetrievedBy:    creator,
	}
	queryHandler := contracttemplate.GetByIdHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	assert.Equal(t, contractTemplate.DID, *did)
	assert.Equal(t, templatestate.Draft, contractTemplate.State)
}

func TestRetrieve_RetrieveAllContractTemplates(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	creator := "Test User"

	dids := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		did, err := base.GetDID()
		if err != nil {
			t.Fatalf("Failed to get new DID: %v", err)
		}
		dids = append(dids, *did)
		createContractTemplate(t, db, did, templatestate.Draft, creator)
	}
	sort.Strings(dids)

	ctx := context.Background()

	qry := contracttemplate.GetAllMetaDataQuery{
		RetrievedBy: creator,
	}
	queryHandler := contracttemplate.GetAllMetaDataHandler{
		Ctx: ctx,
		DB:  db,
	}
	result, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	for _, ct := range result.ContractTemplates {
		assert.Equal(t, templatestate.Draft, ct.State)

		if !slices.Contains(dids, ct.DID) {
			t.Errorf("DID not found in retrieved contract template: %v", ct.DID)
		}
	}
}
