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

func TestSubmit_RetrieveContractTemplateById(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	did, err := base.GetDID()
	if err != nil {
		t.Fatalf("Failed to connect get new DID: %v", err)
	}

	createTestContractTemplate(t, did, templatestate.Draft, db)

	ctx := context.Background()

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

	assert.Equal(t, contractTemplate.DID, *did)
	assert.Equal(t, templatestate.Draft, contractTemplate.State)
}

func TestSubmit_RetrieveAllContractTemplate(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	dids := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		did, err := base.GetDID()
		if err != nil {
			t.Fatalf("Failed to connect get new DID: %v", err)
		}
		dids = append(dids, *did)
		createTestContractTemplate(t, did, templatestate.Draft, db)
	}
	sort.Strings(dids)

	ctx := context.Background()

	retrievedBy := "Test User"

	qry := contracttemplate.GetAllContractTemplatesMetaDataByFilterQuery{
		RetrievedBy: retrievedBy,
	}
	queryHandler := contracttemplate.GetAllContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	for _, ct := range contractTemplate {
		assert.Equal(t, templatestate.Draft, ct.State)

		if !slices.Contains(dids, ct.DID) {
			t.Errorf("DID not found in retrieved contract template: %v", ct.DID)
		}
	}
}
