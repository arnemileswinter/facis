package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	templateevents "digital-contracting-service/internal/templaterepository/event"

	"github.com/jmoiron/sqlx"
)

type CreateTemplateContractCommand struct {
	DID         string
	CreatedBy   string
	Name        *string
	Description *string
	MetaData    *datatype.JSON
}

type CreateTemplateContractHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *CreateTemplateContractHandler) Handle(cmd CreateTemplateContractCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.GetTransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	data := templaterepository.ContractTemplateData{
		DID:         cmd.DID,
		CreatedBy:   cmd.CreatedBy,
		State:       templatestate.Draft,
		Name:        cmd.Name,
		Description: cmd.Description,
		MetaData:    cmd.MetaData,
	}
	createdAt, err := templaterepository.CreateContractTemplate(ctx, tx, data)
	if err != nil {
		return err
	}

	evt := templateevents.ContractTemplateCreatedEvent{
		DID:            cmd.DID,
		DocumentNumber: 1,
		Version:        1,
		CreatedBy:      cmd.CreatedBy,
		Name:           cmd.Name,
		Description:    cmd.Description,
		MetaData:       cmd.MetaData,
		OccurredAt:     *createdAt,
	}
	err = event.CreateNewEvent(ctx, tx, evt)
	if err != nil {
		return err
	}

	return tx.Commit()
}
