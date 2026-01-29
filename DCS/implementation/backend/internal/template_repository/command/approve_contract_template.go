package command

import (
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type ApproveTemplateContractCommand struct {
	DID                          string
	ApprovedBy                   string
	CurrentContractTemplateState template_state.TemplateState
	DecisionNotes                []string
}

type ApproveTemplateContractHandler struct {
	Db     *sqlx.DB
	Logger *log.Logger
}

func (h *ApproveTemplateContractHandler) Handle(cmd ApproveTemplateContractCommand) error {

	if cmd.CurrentContractTemplateState != template_state.Reviewed {
		return errors.New("current template contract state is invalid")
	}
	query := `UPDATE contract_templates SET
        	state = $2
    	WHERE did = $1
`
	state := template_state.Approved

	result, err := h.Db.Exec(query, cmd.DID, state)
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
