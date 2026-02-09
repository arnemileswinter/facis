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
	DID            string                       `db:"did"`
	DocumentNumber int                          `db:"document_number"`
	Version        int                          `db:"version"`
	State          template_state.TemplateState `db:"state"`
	Name           string                       `db:"name"`
	Description    string                       `db:"description"`
	CreatedBy      string                       `db:"created_by"`
	CreatedAt      time.Time                    `db:"created_at"`
	UpdatedBy      string                       `db:"updated_by"`
	UpdatedAt      time.Time                    `db:"updated_at"`
	MetaData       datatype.JSON                `db:"meta_data"`
}

type GetAllContractTemplateHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *GetAllContractTemplateHandler) Handle(query GetAllContractTemplatesQuery) ([]GetAllContractTemplateResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

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
