package pg

import (
	"context"
	"digital-contracting-service/internal/templaterepository/db"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type PostgresApprovalTaskRepo struct {
	Ctx context.Context
}

func (r *PostgresApprovalTaskRepo) Create(tx *sqlx.Tx, data db.ApprovalTaskData) (*time.Time, error) {
	statement := `
        INSERT INTO contract_templates_approval_task (
            did, document_number, version, state, approver, created_by
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING created_at
    `
	var createdAt time.Time
	err := tx.GetContext(r.Ctx, &createdAt, statement,
		data.DID, data.DocumentNumber, data.Version,
		data.State, data.Approver, data.CreatedBy,
	)
	if err != nil {
		return nil, err
	}
	return &createdAt, nil
}

func (r *PostgresApprovalTaskRepo) ReopenTasks(tx *sqlx.Tx, did string, documentNumber int, version int) error {
	statement := `
        UPDATE contract_templates_approval_task SET state = 'OPEN'
        WHERE did = $1 AND document_number = $2 AND version = $3
    `
	_, err := tx.ExecContext(r.Ctx, statement, did, documentNumber, version)
	return err
}

func (r *PostgresApprovalTaskRepo) ReadAll(tx *sqlx.Tx, did string) ([]db.ApprovalTaskData, error) {
	query := `
        SELECT id, did, document_number, version, state, approver,
               created_by, created_at
        FROM contract_templates_approval_task WHERE did = $1
    `
	var approvalTasks []db.ApprovalTaskData
	err := tx.SelectContext(r.Ctx, &approvalTasks, query, did)
	if err != nil {
		return nil, err
	}
	return approvalTasks, nil
}

func (r *PostgresApprovalTaskRepo) ReadAllByApprover(tx *sqlx.Tx, approver string) ([]db.ApprovalTaskData, error) {
	query := `
        SELECT id, did, document_number, version, state, approver,
               created_by, created_at
        FROM contract_templates_approval_task WHERE approver = $1
    `
	var approvalTasks []db.ApprovalTaskData
	err := tx.SelectContext(r.Ctx, &approvalTasks, query, approver)
	if err != nil {
		return nil, err
	}
	return approvalTasks, nil
}

func (r *PostgresApprovalTaskRepo) Update(tx *sqlx.Tx, did string, documentNumber int, version int, approver string, state string) error {
	statement := `
        UPDATE contract_templates_approval_task SET state = $5
        WHERE did = $1 AND document_number = $2 AND version = $3 AND approver = $4
    `
	result, err := tx.ExecContext(r.Ctx, statement, did, documentNumber, version, approver, state)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user has no review task for this contract template")
	}
	return nil
}

func (r *PostgresApprovalTaskRepo) IsValidApprover(tx *sqlx.Tx, did string, documentNumber int, version int, approver string) (bool, error) {
	query := `
        SELECT COUNT(*) FROM contract_templates_approval_task
        WHERE did = $1 AND document_number = $2 AND version = $3 AND approver = $4
    `
	var count int
	err := tx.GetContext(r.Ctx, &count, query, did, documentNumber, version, approver)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresApprovalTaskRepo) TaskExistsInState(tx *sqlx.Tx, did string, documentNumber int, version int, approver string, state string) (bool, error) {
	query := `
        SELECT COUNT(*) FROM contract_templates_approval_task
        WHERE did = $1 AND document_number = $2 AND version = $3 AND approver = $4 AND state = $5
    `
	var count int
	err := tx.GetContext(r.Ctx, &count, query, did, documentNumber, version, approver, state)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresApprovalTaskRepo) TaskExists(tx *sqlx.Tx, did string, documentNumber int, version int) (bool, error) {
	query := `
        SELECT COUNT(*) FROM contract_templates_approval_task
        WHERE did = $1 AND document_number = $2 AND version = $3
    `
	var count int
	err := tx.GetContext(r.Ctx, &count, query, did, documentNumber, version)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresApprovalTaskRepo) Delete(tx *sqlx.Tx, did string, documentNumber int, version int) error {
	statement := `
        DELETE FROM contract_templates_approval_task
        WHERE did = $1 AND document_number = $2 AND version = $3
    `
	_, err := tx.ExecContext(r.Ctx, statement, did, documentNumber, version)
	return err
}
