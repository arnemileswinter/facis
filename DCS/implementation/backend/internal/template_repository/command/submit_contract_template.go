package command

import (
	"context"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/template_repository"
	"digital-contracting-service/internal/template_repository/datatype/action_flag"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	templateevents "digital-contracting-service/internal/template_repository/event"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type SubmitTemplateContractCommand struct {
	DID            string
	SubmittedBy    string
	ActionFlag     *action_flag.ActionFlag
	ReviewComments []string
}

type SubmitTemplateContractHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *SubmitTemplateContractHandler) Handle(cmd SubmitTemplateContractCommand) error {

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

	var nextTemplateState template_state.TemplateState
	if *currentTemplateState == template_state.Draft {

		nextTemplateState = template_state.Submitted

	} else if *currentTemplateState == template_state.Submitted {

		if cmd.ActionFlag != nil {
			if *cmd.ActionFlag == action_flag.Approval {
				nextTemplateState = template_state.Reviewed
			} else if *cmd.ActionFlag == action_flag.Draft {
				nextTemplateState = template_state.Draft
			}
		} else {
			return errors.New("action flags is missing")
		}

	} else if *currentTemplateState == template_state.Reviewed {

		nextTemplateState = template_state.Submitted

	} else {
		return errors.New("current template contract state is invalid")
	}

	query := `UPDATE contract_templates SET
        	state = $2
    	WHERE did = $1
`

	_, err = tx.ExecContext(ctx, query, cmd.DID, nextTemplateState)
	if err != nil {
		return err
	}

	evt := templateevents.ContractTemplateSubmittedEvent{
		DID:            cmd.DID,
		SubmittedBy:    cmd.SubmittedBy,
		PreviousState:  *currentTemplateState,
		NewState:       nextTemplateState,
		ActionFlag:     cmd.ActionFlag,
		ReviewComments: cmd.ReviewComments,
		OccurredAt:     time.Now(),
	}
	err = event.CreateNewEvent(h.Ctx, tx, evt)
	if err != nil {
		return err
	}

	return tx.Commit()
}
