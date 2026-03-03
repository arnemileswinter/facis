package reviewtask

import (
	"context"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtask"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type PostgresReviewTaskRepo struct {
	Ctx context.Context
}

func (r *PostgresReviewTaskRepo) Create(tx *sqlx.Tx, data reviewtask.TaskData) (*time.Time, error) {
	statement := `
        INSERT INTO contract_templates_review_task (
            did, document_number, version, state, reviewer, created_by
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING created_at
    `
	var createdAt time.Time
	err := tx.GetContext(r.Ctx, &createdAt, statement,
		data.DID, data.DocumentNumber, data.Version,
		data.State, data.Reviewer, data.CreatedBy,
	)
	if err != nil {
		return nil, err
	}
	return &createdAt, nil
}

func (r *PostgresReviewTaskRepo) IsValidReviewer(tx *sqlx.Tx, did string, documentNumber int, version int, reviewer string) (bool, error) {
	query := `
        SELECT COUNT(*) FROM contract_templates_review_task
        WHERE did = $1 AND document_number = $2 AND version = $3 AND reviewer = $4
    `
	var count int
	err := tx.GetContext(r.Ctx, &count, query, did, documentNumber, version, reviewer)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresReviewTaskRepo) ReopenTasks(tx *sqlx.Tx, did string, documentNumber int, version int) error {
	statement := `
        UPDATE contract_templates_review_task SET state = 'OPEN'
        WHERE did = $1 AND document_number = $2 AND version = $3
    `
	_, err := tx.ExecContext(r.Ctx, statement, did, documentNumber, version)
	return err
}

func (r *PostgresReviewTaskRepo) ReadAll(tx *sqlx.Tx, did string) ([]reviewtask.TaskData, error) {
	query := `
        SELECT id, did, document_number, version, state, reviewer,
               created_by, created_at
        FROM contract_templates_review_task WHERE did = $1
    `
	var reviewTasks []reviewtask.TaskData
	err := tx.SelectContext(r.Ctx, &reviewTasks, query, did)
	if err != nil {
		return nil, err
	}
	return reviewTasks, nil
}

func (r *PostgresReviewTaskRepo) ReadAllByID(tx *sqlx.Tx, did string, documentNumber int, version int) ([]reviewtask.TaskData, error) {
	query := `
        SELECT id, did, document_number, version, state, reviewer,
               created_by, created_at
        FROM contract_templates_review_task WHERE did = $1 AND document_number = $2 AND version = $3
    `
	var reviewTasks []reviewtask.TaskData
	err := tx.SelectContext(r.Ctx, &reviewTasks, query, did, documentNumber, version)
	if err != nil {
		return nil, err
	}
	return reviewTasks, nil
}

func (r *PostgresReviewTaskRepo) ReadAllByReviewer(tx *sqlx.Tx, reviewer string) ([]reviewtask.TaskData, error) {
	query := `
        SELECT id, did, document_number, version, state, reviewer,
               created_by, created_at
        FROM contract_templates_review_task WHERE reviewer = $1
    `
	var reviewTasks []reviewtask.TaskData
	err := tx.SelectContext(r.Ctx, &reviewTasks, query, reviewer)
	if err != nil {
		return nil, err
	}
	return reviewTasks, nil
}

func (r *PostgresReviewTaskRepo) Update(tx *sqlx.Tx, did string, documentNumber int, version int, reviewer string, state reviewtaskstate.ReviewTaskState) error {
	statement := `
        UPDATE contract_templates_review_task SET state = $5
        WHERE did = $1 AND document_number = $2 AND version = $3 AND reviewer = $4
    `
	result, err := tx.ExecContext(r.Ctx, statement, did, documentNumber, version, reviewer, state)
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

func (r *PostgresReviewTaskRepo) AnyTasksInState(tx *sqlx.Tx, did string, documentNumber int, version int, states ...reviewtaskstate.ReviewTaskState) (bool, error) {
	placeholders := make([]string, len(states))
	args := []interface{}{did, documentNumber, version}

	for i, s := range states {
		placeholders[i] = fmt.Sprintf("$%d", i+4)
		args = append(args, s)
	}

	query := fmt.Sprintf(`
        SELECT COUNT(*) 
        FROM contract_templates_review_task 
        WHERE did = $1 AND document_number = $2 AND version = $3 AND state IN (%s)
    `, strings.Join(placeholders, ", "))

	var count int
	err := tx.GetContext(r.Ctx, &count, query, args...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresReviewTaskRepo) TaskExistsInState(tx *sqlx.Tx, did string, documentNumber int, version int, reviewer string, state reviewtaskstate.ReviewTaskState) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM contract_templates_review_task 
        WHERE did = $1 AND document_number = $2 AND version = $3 AND reviewer = $4 AND state = $5
    `
	var count int
	err := tx.GetContext(r.Ctx, &count, query, did, documentNumber, version, reviewer, state)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresReviewTaskRepo) TaskExist(tx *sqlx.Tx, did string, documentNumber int, version int) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM contract_templates_review_task 
        WHERE did = $1 AND document_number = $2 AND version = $3
    `
	var count int
	err := tx.GetContext(r.Ctx, &count, query, did, documentNumber, version)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresReviewTaskRepo) Delete(tx *sqlx.Tx, did string, documentNumber int, version int) error {
	statement := `
        DELETE FROM contract_templates_review_task
        WHERE did = $1 AND document_number = $2 AND version = $3
    `
	_, err := tx.ExecContext(r.Ctx, statement, did, documentNumber, version)
	return err
}
