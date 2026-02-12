package query

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/template_repository"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	templateevents "digital-contracting-service/internal/template_repository/event"
	"time"

	"github.com/jmoiron/sqlx"
)

type GetContractTemplatesByIdQuery struct {
	DID            string
	DocumentNumber int
	Version        int
	RetrievedBy    string
}

type GetContractTemplatesByIdResult struct {
	DID            string
	DocumentNumber int
	Version        int
	State          template_state.TemplateState
	Name           *string
	Description    *string
	CreatedBy      string
	CreatedAt      time.Time
	MetaData       *datatype.JSON
}

type GetContractTemplateByIdHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *GetContractTemplateByIdHandler) Handle(query GetContractTemplatesByIdQuery) (*GetContractTemplatesByIdResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, base.GetTransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	data, err := template_repository.ReadContractTemplateData(ctx, tx, query.DID, query.DocumentNumber, query.Version)
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

	return &GetContractTemplatesByIdResult{
		DID:            query.DID,
		DocumentNumber: data.DocumentNumber,
		Version:        data.Version,
		State:          data.State,
		Name:           data.Name,
		Description:    data.Description,
		CreatedBy:      data.CreatedBy,
		CreatedAt:      data.CreatedAt,
		MetaData:       data.MetaData,
	}, nil
}
