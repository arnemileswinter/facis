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

type GetContractTemplateByIdQuery struct {
	DID         string
	RetrievedBy string
}

type GetContractTemplateByIdResult struct {
	DID            string                       `db:"did"`
	DocumentNumber int                          `db:"document_number"`
	Version        int                          `db:"version"`
	State          template_state.TemplateState `db:"state"`
	Name           *string                      `db:"name"`
	Description    *string                      `db:"description"`
	CreatedBy      string                       `db:"created_by"`
	CreatedAt      time.Time                    `db:"created_at"`
	UpdatedBy      string                       `db:"updated_by"`
	UpdatedAt      time.Time                    `db:"updated_at"`
	MetaData       *datatype.JSON               `db:"meta_data"`
}

type GetContractTemplateByIdHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *GetContractTemplateByIdHandler) Handle(query GetContractTemplateByIdQuery) (*GetContractTemplateByIdResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	data, err := template_repository.ReadContractTemplateData(ctx, tx, query.DID)
	if err != nil {
		return nil, err
	}

	evt := templateevents.ContractTemplateRetrievedByIdEvent{
		DID:         query.DID,
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

	return &GetContractTemplateByIdResult{
		DID:            query.DID,
		DocumentNumber: data.DocumentNumber,
		Version:        data.Version,
		State:          data.State,
		Name:           data.Name,
		Description:    data.Description,
		CreatedBy:      data.CreatedBy,
		CreatedAt:      data.CreatedAt,
		UpdatedBy:      data.UpdatedBy,
		UpdatedAt:      data.UpdatedAt,
		MetaData:       data.MetaData,
	}, nil
}
