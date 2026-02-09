package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"digital-contracting-service/internal/template_repository/query"
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

	createTestContractTemplate(t, did, template_state.Draft, db)

	ctx := context.Background()

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

	assert.Equal(t, contractTemplate.DID, *did)
	assert.Equal(t, template_state.Draft, contractTemplate.State)
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
		createTestContractTemplate(t, did, template_state.Draft, db)
	}
	sort.Strings(dids)

	ctx := context.Background()

	retrievedBy := "Test User"

	qry := query.GetAllContractTemplatesQuery{
		RetrievedBy: retrievedBy,
	}
	queryHandler := query.GetAllContractTemplateHandler{
		Ctx: ctx,
		DB:  db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		t.Fatalf("Failed to query template contract: %v", err)
	}

	for i, ct := range contractTemplate {
		assert.Equal(t, dids[i], ct.DID)
		assert.Equal(t, template_state.Draft, ct.State)
	}
}
