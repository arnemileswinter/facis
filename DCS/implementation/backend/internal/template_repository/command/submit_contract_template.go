package command

import (
	"digital-contracting-service/internal/template_repository/datatype/action_flag"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type SubmitTemplateContractCommand struct {
	DID                          string
	SubmittedBy                  string
	CurrentContractTemplateState template_state.TemplateState
	ActionFlag                   *action_flag.ActionFlag
	ReviewComments               []string
}

type SubmitTemplateContractHandler struct {
	DB     *sqlx.DB
	Logger *log.Logger
}

func (h *SubmitTemplateContractHandler) Handle(cmd SubmitTemplateContractCommand) error {

	var nextTemplateState template_state.TemplateState
	if cmd.CurrentContractTemplateState == template_state.Draft {

		nextTemplateState = template_state.Submitted

	} else if cmd.CurrentContractTemplateState == template_state.Submitted {

		if cmd.ActionFlag != nil {
			if *cmd.ActionFlag == action_flag.Approval {
				nextTemplateState = template_state.Reviewed
			} else if *cmd.ActionFlag == action_flag.Draft {
				nextTemplateState = template_state.Draft
			}
		} else {
			return errors.New("action flags is missing")
		}

	} else if cmd.CurrentContractTemplateState == template_state.Reviewed {

		if cmd.ActionFlag != nil {
			if *cmd.ActionFlag == action_flag.Draft {
				nextTemplateState = template_state.Draft
			} else {
				return errors.New("invalid action flag for this contract template state")
			}
		} else {
			return errors.New("action flags is missing")
		}

	} else {
		return errors.New("current template contract state is invalid")
	}

	query := `UPDATE contract_templates SET
        	state = $2
    	WHERE did = $1
`

	result, err := h.DB.Exec(query, cmd.DID, nextTemplateState)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("couldn't update contract template state")
	}

	// submittedAt := time.Now()

	return nil
}
