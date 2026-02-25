package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/datatype/templatetype"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"digital-contracting-service/internal/templaterepository/reviewtask"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type UpdateCommand struct {
	DID            string
	DocumentNumber int
	Version        int
	TemplateType   *templatetype.TemplateType
	UpdatedAt      time.Time
	UpdatedBy      string
	Name           *string
	Description    *string
	TemplateData   *datatype.JSON
}

type UpdateHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *UpdateHandler) Handle(cmd UpdateCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	oldData, err := templaterepository.ReadDataByID(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not read template data: %w", err)
	}

	if cmd.UpdatedAt.Before(oldData.UpdatedAt) {
		return errors.New("contract template was updated elsewhere, please reload")
	}

	if oldData.State != templatestate.Draft && oldData.State != templatestate.Submitted {
		return errors.New("invalid contract template state")
	}

	isValidUser := false
	if oldData.State == templatestate.Draft && oldData.CreatedBy == cmd.UpdatedBy {
		isValidUser = true
	} else if oldData.State == templatestate.Submitted {
		valid, err := reviewtask.IsValidTaskUser(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, cmd.UpdatedBy)
		if err != nil {
			return err
		}
		isValidUser = valid
	}

	if !isValidUser {
		return fmt.Errorf("invalid user")
	}

	err = templaterepository.ReopenTasks(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not reopen tasks: %w", err)
	}

	newData := templaterepository.ContractTemplateUpdateData{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		TemplateType:   cmd.TemplateType,
		Name:           cmd.Name,
		Description:    cmd.Description,
		TemplateData:   cmd.TemplateData,
	}
	err = templaterepository.UpdateData(ctx, tx, newData)
	if err != nil {
		return fmt.Errorf("could not update template data: %w", err)
	}

	evt := templateevents.UpdateContractTemplateEvent{
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
