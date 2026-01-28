package command

import (
	"digital-contracting-service/internal/base/datatype"
	"errors"
	"log"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type UpdateTemplateContractCommand struct {
	DID       string
	UpdatedBy string

	Name        *string
	Description *string
	MetaData    *datatype.JSON
}

type UpdateTemplateContractHandler struct {
	Db     *sqlx.DB
	Logger *log.Logger
}

func (h *UpdateTemplateContractHandler) Handle(cmd UpdateTemplateContractCommand) error {
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

	err := h.Db.QueryRow(query, params...).Err()
	if err != nil {
		return err
	}

	// updatedAt := time.Now()

	return nil
}
