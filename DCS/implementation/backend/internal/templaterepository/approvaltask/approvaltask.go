package approvaltask

import (
	"context"
	aopprovaltaskstate "digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskData struct {
	ID             string                               `db:"id"`
	DID            string                               `db:"did"`
	DocumentNumber int                                  `db:"document_number"`
	Version        int                                  `db:"version"`
	State          aopprovaltaskstate.ApprovalTaskState `db:"state"`
	Approver       string                               `db:"approver"`
	CreatedBy      string                               `db:"created_by"`
	CreatedAt      time.Time                            `db:"created_at"`
}

func CreateTask(ctx context.Context, tx *sqlx.Tx, data TaskData) (*time.Time, error) {
	statement := `
    INSERT INTO contract_templates_approval_task (
        did, document_number, version, state, approver, created_by
    ) VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING created_at
`

	var createdAt time.Time
	err := tx.GetContext(ctx, &createdAt, statement,
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

func ReopenTasks(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) error {
	statement := `
        UPDATE contract_templates_approval_task SET state = 'OPEN'
        WHERE did = $1 AND document_number = $2 AND version = $3
    `

	_, err := tx.ExecContext(ctx, statement, did, documentNumber, version)
	if err != nil {
		return err
	}

	return nil
}

func ReadAll(ctx context.Context, tx *sqlx.Tx, did string) ([]TaskData, error) {
	query := `
        SELECT id, did, document_number, version, state, approver,
               created_by, created_at
        FROM contract_templates_approval_task WHERE did = $1
    `

	var approvalTasks []TaskData
	err := tx.SelectContext(ctx, &approvalTasks, query, did)
	if err != nil {
		return nil, err
	}
	return approvalTasks, nil
}

func ReadAllByApprover(ctx context.Context, tx *sqlx.Tx, approver string) ([]TaskData, error) {
	query := `
        SELECT id, did, document_number, version, state, approver,
               created_by, created_at
        FROM contract_templates_approval_task WHERE approver = $1
    `

	var approvalTasks []TaskData
	err := tx.SelectContext(ctx, &approvalTasks, query, approver)
	if err != nil {
		return nil, err
	}
	return approvalTasks, nil
}

func IsValidTaskUser(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int, approver string) (bool, error) {
	selectQuery := `
        SELECT COUNT(*) FROM contract_templates_approval_task
		WHERE did = $1 AND document_number = $2 AND version = $3 AND approver = $4
`

	var count int
	err := tx.GetContext(ctx, &count, selectQuery, did, documentNumber, version, approver)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func DeleteTask(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) error {
	statement := `
        DELETE FROM contract_templates_approval_task
        WHERE did = $1 AND document_number = $2 AND version = $3
    `

	_, err := tx.ExecContext(ctx, statement, did, documentNumber, version)
	if err != nil {
		return err
	}

	return err
}
