package query

import (
	"context"
	"database/sql"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type GetContractTemplateQuery struct {
	DID string
}

type GetContractTemplateResponse struct {
	DID            string                       `db:"did"`
	DocumentNumber int                          `db:"document_number"`
	Version        int                          `db:"version"`
	State          template_state.TemplateState `db:"state"`
	Name           string                       `db:"name"`
	Description    string                       `db:"description"`
	CreatedBy      string                       `db:"created_by"`
	CreatedAt      time.Time                    `db:"created_at"`
	MetaData       datatype.JSON                `db:"meta_data"`
}

type GetContractTemplateHandler struct {
	Ctx    context.Context
	Db     *sqlx.DB
	Logger *log.Logger
}

func (h *GetContractTemplateHandler) GetContractTemplate(query GetContractTemplateQuery) (*GetContractTemplateResponse, error) {
	sqlQuery := `
        SELECT 
            did,
            document_number,
            version,
            state,
            name,
            description,
            created_by,
            created_at,
            meta_data
        FROM contract_templates
        WHERE did = $1
    `

	var contractTemplate GetContractTemplateResponse
	err := h.Db.GetContext(h.Ctx, &contractTemplate, sqlQuery, query.DID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("contract template with DID %s not found", query.DID))
		}
		return nil, err
	}

	return &contractTemplate, nil
}
