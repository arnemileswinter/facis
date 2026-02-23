package contracttemplate

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type GetAllContractTemplatesMetaDataByFilterQuery struct {
	RetrievedBy string

	DID            *string
	DocumentNumber *int
	Version        *int
	State          *templatestate.TemplateState
	Name           *string
	Description    *string
	Filter         *string
}

type GetAllContractTemplatesMetaDataByFilterResult struct {
	DID            string
	DocumentNumber int
	Version        int
	State          templatestate.TemplateState
	Name           string
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	MetaData       datatype.JSON
}

type GetAllContractTemplatesMetaDataByFilterHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *GetAllContractTemplatesMetaDataByFilterHandler) Handle(query GetAllContractTemplatesMetaDataByFilterQuery) ([]GetAllContractTemplatesMetaDataByFilterResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create transaction: %w", err)
	}
	defer tx.Rollback()

	searchValues := templaterepository.SearchValues{
		DID:            query.DID,
		DocumentNumber: query.DocumentNumber,
		Version:        query.Version,
		State:          query.State,
		Name:           query.Name,
		Description:    query.Description,
		Filter:         query.Filter,
	}

	contractTemplates, err := templaterepository.ReadAllMetaDataByFilter(ctx, tx, searchValues)
	if err != nil {
		return nil, fmt.Errorf("could not read all contract templates: %w", err)
	}

	evt := templateevents.RetrieveAllContractTemplatesEvent{
		RetrievedBy: query.RetrievedBy,
		OccurredAt:  time.Now(),
	}
	err = event.Create(h.Ctx, tx, evt)
	if err != nil {
		return nil, fmt.Errorf("could not create event: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	result := make([]GetAllContractTemplatesMetaDataByFilterResult, len(contractTemplates))
	for i, data := range contractTemplates {
		result[i] = GetAllContractTemplatesMetaDataByFilterResult{
			DID:            data.DID,
			DocumentNumber: data.DocumentNumber,
			Version:        data.Version,
			State:          data.State,
			Name:           *data.Name,
			Description:    *data.Description,
			CreatedAt:      data.CreatedAt,
			UpdatedAt:      data.UpdatedAt,
		}
	}

	return result, nil
}
