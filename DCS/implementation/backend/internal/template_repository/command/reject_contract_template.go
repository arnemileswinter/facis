package command

import (
	"context"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	templateevents "digital-contracting-service/internal/template_repository/event"
	"time"

	"github.com/jmoiron/sqlx"
)

type RejectTemplateContractCommand struct {
	DID                          string
	RejectedBy                   string
	CurrentContractTemplateState template_state.TemplateState
	Reason                       string
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

	//currentState, err := template_repository.ReadContractTemplateState(ctx, tx, cmd.DID)
	//if err != nil {
	//	return err
	//}

	//if *currentState != template_state.Draft {
	//	return errors.New("invalid contract template state")
	//}

	query := `UPDATE contract_templates SET
        	state = $2
    	WHERE did = $1
`

	_, err = tx.ExecContext(ctx, query, cmd.DID, template_state.Approved)
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

	return nil

	return nil
}
