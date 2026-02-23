package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CreateCommand struct {
	DID          string
	CreatedBy    string
	Name         *string
	Description  *string
	TemplateData *datatype.JSON
}

type CreateHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *CreateHandler) Handle(cmd CreateCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	data := templaterepository.ContractTemplate{
		DID:          cmd.DID,
		CreatedBy:    cmd.CreatedBy,
		State:        templatestate.Draft,
		Name:         cmd.Name,
		Description:  cmd.Description,
		TemplateData: cmd.TemplateData,
	}
	createdAt, err := templaterepository.CreateContractTemplate(ctx, tx, data)
	if err != nil {
		return fmt.Errorf("could not create contract template: %w", err)
	}

	evt := templateevents.CreateContractTemplateEvent{
		DID:            cmd.DID,
		DocumentNumber: 1,
		Version:        1,
		CreatedBy:      cmd.CreatedBy,
		Name:           cmd.Name,
		Description:    cmd.Description,
		TemplateData:   cmd.TemplateData,
		OccurredAt:     *createdAt,
	}
	err = event.Create(ctx, tx, evt)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	return tx.Commit()
}
