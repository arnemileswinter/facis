package template_repository

import (
	"context"
	"digital-contracting-service/internal/template_repository/datatype/review_task_state"
	"time"

	"github.com/jmoiron/sqlx"
)

type ReviewTaskData struct {
	ID             string                            `db:"id"`
	DID            string                            `db:"did"`
	DocumentNumber int                               `db:"document_number"`
	Version        int                               `db:"version"`
	State          review_task_state.ReviewTaskState `db:"state"`
	Reviewer       string                            `db:"reviewer"`
	CreatedBy      string                            `db:"created_by"`
	CreatedAt      time.Time                         `db:"created_at"`
}

func CreateReviewTasks(ctx context.Context, tx *sqlx.Tx, data ReviewTaskData) (*time.Time, error) {
	query := `
    INSERT INTO contract_templates_review_task (
        did, document_number, version, state, reviewer, created_by
    ) VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING created_at
`

	var createdAt time.Time
	err := tx.GetContext(ctx, &createdAt, query,
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

func ReopenReviewTasks(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) error {
	query := `
        UPDATE contract_templates_review_task SET state = 'OPEN'
        WHERE did = $1 AND document_number = $2 AND version = $3
    `

	_, err := tx.ExecContext(ctx, query, did, documentNumber, version)
	if err != nil {
		return err
	}

	return err
}

func ReadAllReviewTasks(ctx context.Context, tx *sqlx.Tx, did string) ([]ReviewTaskData, error) {
	query := `
        SELECT id, did, document_number, version, state, reviewer,
               created_by, created_at
        FROM contract_templates_review_task WHERE did = $1
    `

	var reviewTasks []ReviewTaskData
	err := tx.SelectContext(ctx, &reviewTasks, query, did)
	if err != nil {
		return nil, err
	}
	return reviewTasks, nil
}

func UpdateReviewTask(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int, reviewer string, state review_task_state.ReviewTaskState) error {
	query := `
        UPDATE contract_templates_review_task SET state = $5
        WHERE did = $1 AND document_number = $2 AND version = $3 AND reviewer = $4
    `

	_, err := tx.ExecContext(ctx, query, did, documentNumber, version, reviewer, state)
	if err != nil {
		return err
	}

	return err
}

func ExistReviewTaskInState(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int, state review_task_state.ReviewTaskState) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM contract_templates_review_task 
        WHERE did = $1 AND document_number = $2 AND version = $3 AND state = $4
    `

	var count int
	err := tx.GetContext(ctx, &count, query, did, documentNumber, version, state)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func DeleteReviewTask(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) error {
	query := `
        DELETE FROM contract_templates_review_task
        WHERE did = $1 AND document_number = $2 AND version = $3
    `

	_, err := tx.ExecContext(ctx, query, did, documentNumber, version)
	if err != nil {
		return err
	}

	return err
}
