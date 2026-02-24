package reviewtask

import (
	"context"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type TaskData struct {
	ID             string                          `db:"id"`
	DID            string                          `db:"did"`
	DocumentNumber int                             `db:"document_number"`
	Version        int                             `db:"version"`
	State          reviewtaskstate.ReviewTaskState `db:"state"`
	Reviewer       string                          `db:"reviewer"`
	CreatedBy      string                          `db:"created_by"`
	CreatedAt      time.Time                       `db:"created_at"`
}

func CreateTask(ctx context.Context, tx *sqlx.Tx, data TaskData) (*time.Time, error) {
	statement := `
    INSERT INTO contract_templates_review_task (
        did, document_number, version, state, reviewer, created_by
    ) VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING created_at
`

	var createdAt time.Time
	err := tx.GetContext(ctx, &createdAt, statement,
		data.DID,
		data.DocumentNumber,
		data.Version,
		data.State,
		data.Reviewer,
		data.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &createdAt, nil
}

func IsValidTaskUser(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int, reviewer string) (bool, error) {
	query := `
        SELECT COUNT(*) FROM contract_templates_review_task
		WHERE did = $1 AND document_number = $2 AND version = $3 AND reviewer = $4
`

	var count int
	err := tx.GetContext(ctx, &count, query, did, documentNumber, version, reviewer)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func ReopenTasks(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) error {
	statement := `
        UPDATE contract_templates_review_task SET state = 'OPEN'
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
        SELECT id, did, document_number, version, state, reviewer,
               created_by, created_at
        FROM contract_templates_review_task WHERE did = $1
    `

	var reviewTasks []TaskData
	err := tx.SelectContext(ctx, &reviewTasks, query, did)
	if err != nil {
		return nil, err
	}
	return reviewTasks, nil
}

func ReadAllByReviewer(ctx context.Context, tx *sqlx.Tx, reviewer string) ([]TaskData, error) {
	query := `
        SELECT id, did, document_number, version, state, reviewer,
               created_by, created_at
        FROM contract_templates_review_task WHERE reviewer = $1
    `

	var reviewTasks []TaskData
	err := tx.SelectContext(ctx, &reviewTasks, query, reviewer)
	if err != nil {
		return nil, err
	}
	return reviewTasks, nil
}

func UpdateTask(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int, reviewer string, state reviewtaskstate.ReviewTaskState) error {
	statement := `
        UPDATE contract_templates_review_task SET state = $5
        WHERE did = $1 AND document_number = $2 AND version = $3 AND reviewer = $4
    `

	result, err := tx.ExecContext(ctx, statement, did, documentNumber, version, reviewer, state)
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

	return err
}

func ExistTasksInStates(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int, states ...reviewtaskstate.ReviewTaskState) (bool, error) {
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
	err := tx.GetContext(ctx, &count, query, args...)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func HasTaskInState(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int, reviewer string, state reviewtaskstate.ReviewTaskState) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM contract_templates_review_task 
        WHERE did = $1 AND document_number = $2 AND version = $3 AND reviewer = $4 AND state = $5
    `

	var count int
	err := tx.GetContext(ctx, &count, query, did, documentNumber, version, reviewer, state)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func DeleteTask(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) error {
	statement := `
        DELETE FROM contract_templates_review_task
        WHERE did = $1 AND document_number = $2 AND version = $3
    `

	_, err := tx.ExecContext(ctx, statement, did, documentNumber, version)
	if err != nil {
		return err
	}

	return err
}
