package query

import (
	"context"
	"database/sql"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type GetAllContractTemplateReviewTasksForDID struct {
	DID string
}

type GetAllContractTemplateReviewTasksForDIDResult struct {
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

type GetAllContractTemplateReviewTasksForDIDHandler struct {
	Ctx context.Context
	Db  *sqlx.DB
}

func (h *GetAllContractTemplateReviewTasksForDIDHandler) Handle(query GetAllContractTemplateReviewTasksForDID) (*GetAllContractTemplateReviewTasksForDIDResult, error) {
	sqlQuery := `
        SELECT 
            did,
            document_number,
            version,
            state,
            created_by,
            created_at,
            updated_by,
            updated_at
        FROM contract_templates_review_task
        WHERE did = $1
    `

	var result GetAllContractTemplateReviewTasksForDIDResult
	err := h.Db.GetContext(h.Ctx, &result, sqlQuery, query.DID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("contract template with DID %s not found", query.DID))
		}
		return nil, err
	}

	return &result, nil
}
