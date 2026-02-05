package command

import (
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"errors"
	"log"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type UpdateTemplateContractCommand struct {
	DID       string
	UpdatedBy string

	CurrentContractTemplateState template_state.TemplateState

	Name        *string
	Description *string
	MetaData    *datatype.JSON
}

type UpdateTemplateContractHandler struct {
	Db     *sqlx.DB
	Logger *log.Logger
}

func (h *UpdateTemplateContractHandler) Handle(cmd UpdateTemplateContractCommand) error {

	if cmd.CurrentContractTemplateState != template_state.Draft {
		return errors.New("invalid contract template state")
	}

	query := `UPDATE contract_templates SET`

	var params []interface{}
	paramIndex := 1

	if cmd.Name != nil {
		query += ` name = $` + strconv.Itoa(paramIndex) + `,`
		params = append(params, cmd.Name)
		paramIndex++
	}

	if cmd.Description != nil {
		query += ` description = $` + strconv.Itoa(paramIndex) + `,`
		params = append(params, cmd.Description)
		paramIndex++
	}

	if cmd.MetaData != nil && cmd.MetaData.IsNotNullValue() {
		query += ` meta_data = $` + strconv.Itoa(paramIndex) + `,`
		params = append(params, cmd.MetaData)
		paramIndex++
	}

	if len(params) == 0 {
		return errors.New("no parameters found")
	}

	// Remove last comma
	query = query[:len(query)-1]

	query += ` WHERE did = $` + strconv.Itoa(paramIndex) + `;`
	params = append(params, cmd.DID)

	result, err := h.Db.Exec(query, params...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("couldn't update template contract")
	}

	// updatedAt := time.Now()

	return nil
}
