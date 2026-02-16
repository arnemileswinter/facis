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
	Filter      map[string]interface{}
}

type GetAllContractTemplatesMetaDataResult struct {
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

type GetAllContractTemplateHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *GetAllContractTemplateHandler) Handle(query GetAllContractTemplatesMetaDataByFilterQuery) ([]GetAllContractTemplatesMetaDataResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create transaction: %w", err)
	}
	defer tx.Rollback()

	contractTemplates, err := templaterepository.ReadAllContractTemplateMetaDataByFilter(ctx, tx, query.Filter)
	if err != nil {
		return nil, fmt.Errorf("could not read all contract templates: %w", err)
	}

	evt := templateevents.ContractTemplateRetrievedAllEvent{
		RetrievedBy: query.RetrievedBy,
		Filter:      query.Filter,
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

	result := make([]GetAllContractTemplatesMetaDataResult, len(contractTemplates))
	for i, data := range contractTemplates {
		result[i] = GetAllContractTemplatesMetaDataResult{
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
