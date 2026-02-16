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
	State          templatestate.TemplateState
	Name           *string
	Description    *string
	CreatedBy      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	TemplateData   *datatype.JSON
}

type GetContractTemplateByIdHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *GetContractTemplateByIdHandler) Handle(query GetContractTemplatesByIdQuery) (*GetContractTemplatesByIdResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	data, err := templaterepository.ReadContractTemplateDataById(ctx, tx, query.DID, query.DocumentNumber, query.Version)
	if err != nil {
		return nil, fmt.Errorf("could not get contract template data: %w", err)
	}

	evt := templateevents.ContractTemplateRetrievedByIdEvent{
		DID:            query.DID,
		DocumentNumber: query.DocumentNumber,
		Version:        query.Version,
		RetrievedBy:    query.RetrievedBy,
		OccurredAt:     time.Now(),
	}
	err = event.Create(h.Ctx, tx, evt)
	if err != nil {
		return nil, fmt.Errorf("could not create event: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
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
		UpdatedAt:      data.UpdatedAt,
		TemplateData:   data.TemplateData,
	}, nil
}
