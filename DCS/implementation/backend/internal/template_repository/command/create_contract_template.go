package command

import (
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type CreateTemplateContractCommand struct {
	DID       string
	CreatedBy string

	Name        *string
	Description *string
	MetaData    *datatype.JSON
}

type CreateTemplateContractHandler struct {
	Db     *sqlx.DB
	Logger *log.Logger
}

func (h *CreateTemplateContractHandler) Handle(cmd CreateTemplateContractCommand) error {
	query := `
    INSERT INTO contract_templates (
        did, created_by, document_number, version, state, name, 
        description, meta_data
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING created_at
`
	documentNumber := 0
	version := 0
	state := template_state.Draft

	var createdAt time.Time
	err := h.Db.QueryRow(
		query,
		cmd.DID,
		cmd.CreatedBy,
		documentNumber,
		version,
		state,
		cmd.Name,
		cmd.Description,
		cmd.MetaData,
	).Scan(&createdAt)

	return err
}
