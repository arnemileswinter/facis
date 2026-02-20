package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type UpdateTemplateContractCommand struct {
	DID            string
	DocumentNumber int
	Version        int
	UpdatedBy      string
	UpdatedAt      time.Time
	Name           *string
	Description    *string
	TemplateData   *datatype.JSON
}

type UpdateTemplateContractHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *UpdateTemplateContractHandler) Handle(cmd UpdateTemplateContractCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	oldData, err := templaterepository.ReadContractTemplateDataById(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not read template data: %w", err)
	}

	if oldData.CreatedBy != cmd.UpdatedBy {
		return fmt.Errorf("invalid user")
	}

	if cmd.UpdatedAt.Before(oldData.UpdatedAt) {
		return errors.New("contract template was updated elsewhere, please reload")
	}

	if oldData.State != templatestate.Draft {
		return errors.New("invalid contract template state")
	}

	newData := templaterepository.ContractTemplateData{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		Name:           cmd.Name,
		Description:    cmd.Description,
		TemplateData:   cmd.TemplateData,
	}
	err = templaterepository.UpdateTemplateContractData(ctx, tx, newData)
	if err != nil {
		return fmt.Errorf("could not update template data: %w", err)
	}

	evt := templateevents.ContractTemplateUpdatedEvent{
		DID:             cmd.DID,
		DocumentNumber:  cmd.DocumentNumber,
		Version:         cmd.Version,
		OldName:         oldData.Name,
		NewName:         cmd.Name,
		OldDescription:  oldData.Description,
		NewDescription:  cmd.Description,
		OldTemplateData: oldData.TemplateData,
		NewTemplateData: cmd.TemplateData,
		UpdatedBy:       cmd.UpdatedBy,
		OccurredAt:      time.Now(),
	}
	err = event.Create(ctx, tx, evt)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	return tx.Commit()
}
