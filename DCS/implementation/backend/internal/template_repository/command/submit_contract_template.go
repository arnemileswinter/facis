package command

import (
	"digital-contracting-service/internal/template_repository/datatype/action_flag"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type SubmitTemplateContractCommand struct {
	DID                   string
	SubmittedBy           string
	ContractTemplateState template_state.TemplateState
	ActionFlag            *action_flag.ActionFlag
	ReviewComments        []string
}

type SubmitTemplateContractHandler struct {
	Db     *sqlx.DB
	Logger *log.Logger
}

func (h *SubmitTemplateContractHandler) Handle(cmd SubmitTemplateContractCommand) error {

	query := `UPDATE contract_templates SET
        	state = $2
    	WHERE did = $1
`

	var state template_state.TemplateState
	if cmd.ContractTemplateState == template_state.Draft {
		state = template_state.Submitted
	} else if cmd.ContractTemplateState == template_state.Submitted {

		if cmd.ActionFlag != nil {
			if *cmd.ActionFlag == action_flag.Approval {
				state = template_state.Approved
			} else if *cmd.ActionFlag == action_flag.Draft {
				state = template_state.Draft
			}
		} else {
			return errors.New("action flags is missing")
		}

	} else {
		return errors.New("invalid state")
	}

	result, err := h.Db.Exec(query, cmd.DID, state)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("couldn't submit contract template")
	}

	// submittedAt := time.Now()

	return nil
}
