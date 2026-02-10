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
	Assignee       string                            `db:"assignee"`
	CreatedBy      string                            `db:"created_by"`
	CreatedAt      time.Time                         `db:"created_at"`
	UpdatedBy      string                            `db:"updated_by"`
	UpdatedAt      time.Time                         `db:"updated_at"`
}

func CreateReviewTask(ctx context.Context, tx *sqlx.Tx, data ReviewTaskData) (*time.Time, error) {
	query := `
    INSERT INTO contract_templates_review_task (
        did, document_number, version, state, assignee, created_by, updated_by
    ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING created_at
`

	var createdAt time.Time
	err := tx.GetContext(ctx, &createdAt, query,
		data.DID,
		data.DocumentNumber,
		data.Version,
		data.State,
		data.Assignee,
		data.CreatedBy,
		data.CreatedBy, // Use created_by for updated_by
	)
	if err != nil {
		return nil, err
	}

	return &createdAt, nil
}

func ReadAllReviewTasks(ctx context.Context, tx *sqlx.Tx, did string) ([]ReviewTaskData, error) {
	query := `
        SELECT id, did, document_number, version, state, assignee,
               created_by, created_at, updated_by, updated_at
        FROM contract_templates_review_task WHERE did = $1
    `

	var reviewTasks []ReviewTaskData
	err := tx.SelectContext(ctx, &reviewTasks, query, did)
	if err != nil {
		return nil, err
	}
	return reviewTasks, nil
}

func UpdateReviewTask(ctx context.Context, tx *sqlx.Tx, did string, version int, assignee string, state review_task_state.ReviewTaskState, reviewComments []string) error {
	query := `
        UPDATE contract_templates SET state = $4, review_comments = $5
        FROM contract_templates_review_task WHERE did = $1 AND version = $2 AND assignee = $3
    `

	var comments string
	for _, comment := range reviewComments {
		comments += comment + ";"
	}
	_, err := tx.ExecContext(ctx, query, did, version, assignee, state, comments)
	if err != nil {
		return err
	}

	return err
}

func CloseReviewTasks(ctx context.Context, tx *sqlx.Tx, did string, version int, assignee string) error {
	query := `
        UPDATE contract_templates SET closed = true
        FROM contract_templates_review_task WHERE did = $1 AND version = $2
    `

	_, err := tx.ExecContext(ctx, query, did, version, assignee)
	if err != nil {
		return err
	}

	return err
}

func ExistReviewTaskInState(ctx context.Context, tx *sqlx.Tx, did string, version int, assignee string, state review_task_state.ReviewTaskState) (bool, error) {
	query := `
        SELECT did, version, state, assignee, state
        FROM contract_templates_review_task WHERE did = $1 AND version = $2 AND assignee = $3 AND state = $4
    `

	var reviewTasks []ReviewTaskData
	err := tx.SelectContext(ctx, &reviewTasks, query, did, version, assignee, state)
	if err != nil {
		return false, err
	}
	return len(reviewTasks) > 0, nil
}
