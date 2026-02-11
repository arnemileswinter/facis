package template_repository

import (
	"context"
	aopprovaltaskstate "digital-contracting-service/internal/template_repository/datatype/approval_task_state"
	"time"

	"github.com/jmoiron/sqlx"
)

type ApprovalTaskData struct {
	ID             string                               `db:"id"`
	DID            string                               `db:"did"`
	DocumentNumber int                                  `db:"document_number"`
	Version        int                                  `db:"version"`
	State          aopprovaltaskstate.ApprovalTaskState `db:"state"`
	Approver       string                               `db:"approver"`
	CreatedBy      string                               `db:"created_by"`
	CreatedAt      time.Time                            `db:"created_at"`
	CancelledAt    *time.Time                           `db:"cancelled_at"`
}

func CreateApprovalTask(ctx context.Context, tx *sqlx.Tx, data ApprovalTaskData) (*time.Time, error) {
	query := `
    INSERT INTO contract_templates_approval_task (
        did, document_number, version, state, approver, created_by
    ) VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING created_at
`

	var createdAt time.Time
	err := tx.GetContext(ctx, &createdAt, query,
		data.DID,
		data.DocumentNumber,
		data.Version,
		data.State,
		data.Approver,
		data.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &createdAt, nil
}

func ReadAllApprovalTasks(ctx context.Context, tx *sqlx.Tx, did string) ([]ApprovalTaskData, error) {
	query := `
        SELECT id, did, document_number, version, state, approver,
               created_by, created_at
        FROM contract_templates_approval_task WHERE did = $1
    `

	var approvalTasks []ApprovalTaskData
	err := tx.SelectContext(ctx, &approvalTasks, query, did)
	if err != nil {
		return nil, err
	}
	return approvalTasks, nil
}

func UpdateApprovalTask(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int, approver string, state aopprovaltaskstate.ApprovalTaskState, comments []string) error {
	query := `
        UPDATE contract_templates_approval_task SET state = $5, comments = $6
        WHERE did = $1 AND document_number = $2 AND version = $3 AND approver = $4
    `

	var commentString string
	for _, comment := range comments {
		commentString += comment + ";"
	}
	_, err := tx.ExecContext(ctx, query, did, documentNumber, version, approver, state, commentString)
	if err != nil {
		return err
	}

	return err
}
