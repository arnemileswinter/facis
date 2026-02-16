package contracttemplate

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"time"

	"github.com/jmoiron/sqlx"
)

type GetAllContractTemplatesQuery struct {
	RetrievedBy string
}

type GetAllContractTemplatesResult struct {
	DID            string
	DocumentNumber int
	Version        int
	State          templatestate.TemplateState
	Name           string
	Description    string
	CreatedBy      string
	CreatedAt      time.Time
	MetaData       datatype.JSON
}

type GetAllContractTemplateHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *GetAllContractTemplateHandler) Handle(query GetAllContractTemplatesQuery) ([]GetAllContractTemplatesResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, base.GetTransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	contractTemplates, err := templaterepository.ReadAllContractTemplateData(ctx, tx)
	if err != nil {
		return nil, err
	}

	evt := templateevents.ContractTemplateRetrievedAllEvent{
		RetrievedBy: query.RetrievedBy,
		OccurredAt:  time.Now(),
	}
	err = event.CreateNewEvent(h.Ctx, tx, evt)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	result := make([]GetAllContractTemplatesResult, len(contractTemplates))
	for i, data := range contractTemplates {
		result[i] = GetAllContractTemplatesResult{
			DID:            data.DID,
			DocumentNumber: data.DocumentNumber,
			Version:        data.Version,
			State:          data.State,
			Name:           *data.Name,
			Description:    *data.Description,
			CreatedBy:      data.CreatedBy,
			CreatedAt:      data.CreatedAt,
			MetaData:       *data.MetaData,
		}
	}

	return result, nil
}
