package command

import (
	"context"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	templateevents "digital-contracting-service/internal/template_repository/event"
	"time"

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

	ctx, cancel := context.WithTimeout(h.Ctx, 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}

	query := `
    INSERT INTO contract_templates (
        did, created_by, updated_by, state, name, 
        description, meta_data
    ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING created_at
`
	state := template_state.Draft

	var createdAt time.Time
	err = tx.QueryRowContext(h.Ctx, query,
		cmd.DID,
		cmd.CreatedBy,
		cmd.CreatedBy, // Set updated_by
		state,
		cmd.Name,
		cmd.Description,
		cmd.MetaData,
	).Scan(&createdAt)
	if err != nil {
		return err
	}

	evt := templateevents.ContractTemplateCreatedEvent{
		DID:         cmd.DID,
		CreatedBy:   cmd.CreatedBy,
		Name:        cmd.Name,
		Description: cmd.Description,
		OccurredAt:  time.Now(),
	}
	err = event.CreateNewEvent(h.Ctx, tx, evt)
	if err != nil {
		return err
	}

	return tx.Commit()
}
