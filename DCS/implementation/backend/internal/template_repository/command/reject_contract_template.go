package command

import (
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type RejectTemplateContractCommand struct {
	DID                          string
	RejectedBy                   string
	CurrentContractTemplateState template_state.TemplateState
	Reason                       string
}

type RejectTemplateContractHandler struct {
	Db     *sqlx.DB
	Logger *log.Logger
}

func (h *RejectTemplateContractHandler) Handle(cmd RejectTemplateContractCommand) error {

	if cmd.CurrentContractTemplateState != template_state.Reviewed {
		return errors.New("current template contract state is invalid")
	}

	if len(cmd.Reason) == 0 {
		return errors.New("reason for rejecting contract template is required")
	}

	query := `UPDATE contract_templates SET
        	state = $2
    	WHERE did = $1
`
	state := template_state.Draft

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
