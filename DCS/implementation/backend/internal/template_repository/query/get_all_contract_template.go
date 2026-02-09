package query

import (
	"context"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type GetAllContractTemplatesQuery struct {
	RetrievedBy string
}

type GetAllContractTemplateResult struct {
	DID            string                       `db:"did"`
	DocumentNumber int                          `db:"document_number"`
	Version        int                          `db:"version"`
	State          template_state.TemplateState `db:"state"`
	Name           string                       `db:"name"`
	Description    string                       `db:"description"`
	CreatedBy      string                       `db:"created_by"`
	CreatedAt      time.Time                    `db:"created_at"`
	UpdatedBy      string                       `db:"updated_by"`
	UpdatedAt      time.Time                    `db:"updated_at"`
	MetaData       datatype.JSON                `db:"meta_data"`
}

type GetAllContractTemplateHandler struct {
	Ctx    context.Context
	Db     *sqlx.DB
	Logger *log.Logger
}

func (h *GetAllContractTemplateHandler) Handle(query GetAllContractTemplatesQuery) ([]GetAllContractTemplateResult, error) {
	sqlQuery := `
        SELECT *
        FROM contract_templates
        ORDER BY did;
    `

	var contractTemplates []GetAllContractTemplateResult
	err := h.Db.SelectContext(h.Ctx, &contractTemplates, sqlQuery)
	if err != nil {
		return nil, err
	}

	return contractTemplates, nil
}
