package command

import (
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type CreateTemplateContractCommand struct {
	DID            string                       `db:"did"`
	DocumentNumber int                          `db:"document_number"`
	Version        int                          `db:"version"`
	State          template_state.TemplateState `db:"state"`
	Name           *string                      `db:"name"`
	Description    *string                      `db:"description"`
	CreatedBy      string                       `db:"created_by"`
	MetaData       *datatype.JSON               `db:"meta_data"`
}

type CreateTemplateContractHandler struct {
	Db     *sqlx.DB
	Logger *log.Logger
}

func (h *CreateTemplateContractHandler) Handle(cmd CreateTemplateContractCommand) error {
	query := `
    INSERT INTO contract_templates (
        did, document_number, version, state, name, 
        description, created_by, meta_data
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING created_at
`
	var createdAt time.Time
	err := h.Db.QueryRow(
		query,
		cmd.DID,
		cmd.DocumentNumber,
		cmd.Version,
		cmd.State,
		cmd.Name,
		cmd.Description,
		cmd.CreatedBy,
		cmd.MetaData,
	).Scan(&createdAt)

	return err
}
