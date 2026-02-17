package test

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
	"fmt"
	"sort"
	"strconv"
	"testing"
)

func TestSearch_SearchContractTemplatesByFilter(t *testing.T) {

	db := setupTestDB(t)

	cleanupContractTemplateTable(t, db)

	states := []templatestate.TemplateState{
		templatestate.Draft,
		templatestate.Submitted,
		templatestate.Rejected,
		templatestate.Reviewed,
		templatestate.Approved,
	}

	dids := make([]string, 0, 50)
	for i := 0; i < 10; i++ {
		did, err := base.GetDID()
		if err != nil {
			t.Fatalf("Failed to connect get new DID: %v", err)
		}
		dids = append(dids, *did)

		templateData := map[string]interface{}{
			"did":            *did,
			"documentNumber": 1 % 5,
			"version":        1 % 5,
			"description":    "test description " + strconv.Itoa(i%5),
		}

		stateId := i % len(states)
		createTestContractTemplateWithTemplateData(t, did, states[stateId], db, templateData)
	}
	sort.Strings(dids)

	ctx := context.Background()

	retrievedBy := "Test User"

	//filter := map[string]interface{}{}

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

	fmt.Println(len(contractTemplate))
}
