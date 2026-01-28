package query

import (
	"context"
	"database/sql"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type GetContractTemplateStateQuery struct {
	DID         string
	RetrievedBy string
}

type GetContractTemplateStateResult struct {
	State template_state.TemplateState `db:"state"`
}

type GetContractTemplateStateHandler struct {
	Ctx    context.Context
	Db     *sqlx.DB
	Logger *log.Logger
}

func (h *GetContractTemplateStateHandler) Handle(query GetContractTemplateStateQuery) (*GetContractTemplateStateResult, error) {
	sqlQuery := `
        SELECT state
        FROM contract_templates
        WHERE did = $1
    `

	var contractTemplate GetContractTemplateStateResult
	err := h.Db.GetContext(h.Ctx, &contractTemplate, sqlQuery, query.DID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contract template with DID %s not found", query.DID)
		}
		return nil, err
	}

	return &contractTemplate, nil
}
