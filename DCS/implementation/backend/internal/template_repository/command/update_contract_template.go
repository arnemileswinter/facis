package command

import (
	"context"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/template_repository"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	templateevents "digital-contracting-service/internal/template_repository/event"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type UpdateTemplateContractCommand struct {
	DID         string
	UpdatedBy   string
	Name        *string
	Description *string
	MetaData    *datatype.JSON
}

type UpdateTemplateContractHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *UpdateTemplateContractHandler) Handle(cmd UpdateTemplateContractCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	oldData, err := template_repository.ReadContractTemplateData(ctx, tx, cmd.DID)
	if err != nil {
		return err
	}

	if oldData.State != template_state.Draft {
		return errors.New("invalid contract template state")
	}

	newData := template_repository.ContractTemplateData{
		DID:         cmd.DID,
		Name:        cmd.Name,
		Description: cmd.Description,
		MetaData:    cmd.MetaData,
	}
	err = template_repository.UpdateTemplateContractData(ctx, tx, newData)
	if err != nil {
		return err
	}

	evt := templateevents.ContractTemplateUpdatedEvent{
		DID:            cmd.DID,
		UpdatedBy:      cmd.UpdatedBy,
		OldName:        oldData.Name,
		NewName:        cmd.Name,
		OldDescription: oldData.Description,
		NewDescription: cmd.Description,
		OldMetaData:    oldData.MetaData,
		NewMetaData:    cmd.MetaData,
		OccurredAt:     time.Now(),
	}
	err = event.CreateNewEvent(ctx, tx, evt)
	if err != nil {
		return err
	}

	return tx.Commit()
}
