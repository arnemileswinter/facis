package command

import (
	"context"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	templateevents "digital-contracting-service/internal/template_repository/event"
	"time"

	"github.com/jmoiron/sqlx"
)

type ApproveTemplateContractCommand struct {
	DID                          string
	ApprovedBy                   string
	CurrentContractTemplateState template_state.TemplateState
	DecisionNotes                []string
}

type ApproveTemplateContractHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *ApproveTemplateContractHandler) Handle(cmd ApproveTemplateContractCommand) error {

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

	evt := templateevents.ContractTemplateApprovedEvent{
		DID:           cmd.DID,
		ApprovedBy:    cmd.ApprovedBy,
		DecisionNotes: cmd.DecisionNotes,
		OccurredAt:    time.Now(),
	}
	err = event.CreateNewEvent(h.Ctx, tx, evt)
	if err != nil {
		return err
	}

	return nil
}
