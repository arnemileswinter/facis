package command

import (
	"context"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/template_repository"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	templateevents "digital-contracting-service/internal/template_repository/event"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type RejectTemplateContractCommand struct {
	DID        string
	RejectedBy string
	Reason     string
}

type RejectTemplateContractHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *RejectTemplateContractHandler) Handle(cmd RejectTemplateContractCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}

	currentTemplateState, err := template_repository.ReadContractTemplateState(ctx, tx, cmd.DID)
	if err != nil {
		return err
	}

	if *currentTemplateState != template_state.Reviewed {
		return errors.New("invalid contract template state")
	}

	query := `UPDATE contract_templates SET
        	state = $2
    	WHERE did = $1
`

	_, err = tx.ExecContext(ctx, query, cmd.DID, template_state.Draft)
	if err != nil {
		return err
	}

	evt := templateevents.ContractTemplateRejectedEvent{
		DID:        cmd.DID,
		RejectedBy: cmd.RejectedBy,
		Reason:     cmd.Reason,
		OccurredAt: time.Now(),
	}
	err = event.CreateNewEvent(h.Ctx, tx, evt)
	if err != nil {
		return err
	}

	return tx.Commit()
}
