package query

import (
	"context"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/template_repository"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	templateevents "digital-contracting-service/internal/template_repository/event"
	"time"

	"github.com/jmoiron/sqlx"
)

type GetAllContractTemplatesQuery struct {
	RetrievedBy string
}

type GetAllContractTemplateResult struct {
	DID            string
	DocumentNumber int
	Version        int
	State          template_state.TemplateState
	Name           string
	Description    string
	CreatedBy      string
	CreatedAt      time.Time
	UpdatedBy      string
	UpdatedAt      time.Time
	MetaData       datatype.JSON
}

type GetAllContractTemplateHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *GetAllContractTemplateHandler) Handle(query GetAllContractTemplatesQuery) ([]GetAllContractTemplateResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	contractTemplates, err := template_repository.ReadAllContractTemplateData(ctx, tx)
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

	result := make([]GetAllContractTemplateResult, len(contractTemplates))
	for i, data := range contractTemplates {
		result[i] = GetAllContractTemplateResult{
			DID:            data.DID,
			DocumentNumber: data.DocumentNumber,
			Version:        data.Version,
			State:          data.State,
			Name:           *data.Name,
			Description:    *data.Description,
			CreatedBy:      data.CreatedBy,
			CreatedAt:      data.CreatedAt,
			UpdatedBy:      data.UpdatedBy,
			UpdatedAt:      data.UpdatedAt,
			MetaData:       *data.MetaData,
		}
	}

	return result, nil
}
