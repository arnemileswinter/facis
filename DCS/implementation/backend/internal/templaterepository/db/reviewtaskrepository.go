package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type ReviewTaskData struct {
	ID             string    `db:"id"`
	DID            string    `db:"did"`
	DocumentNumber string    `db:"document_number"`
	Version        int       `db:"version"`
	State          string    `db:"state"`
	Reviewer       string    `db:"reviewer"`
	CreatedBy      string    `db:"created_by"`
	CreatedAt      time.Time `db:"created_at"`
}

type ReviewTaskRepo interface {
	Create(tx *sqlx.Tx, data ReviewTaskData) (*time.Time, error)
	IsValidReviewer(tx *sqlx.Tx, did string, DocumentNumber string, version int, reviewer string) (bool, error)
	ReopenTasks(tx *sqlx.Tx, did string, DocumentNumber string, version int) error
	ReadAll(tx *sqlx.Tx, did string) ([]ReviewTaskData, error)
	ReadAllByID(tx *sqlx.Tx, did string, DocumentNumber string, version int) ([]ReviewTaskData, error)
	ReadAllByReviewer(tx *sqlx.Tx, reviewer string) ([]ReviewTaskData, error)
	Update(tx *sqlx.Tx, did string, DocumentNumber string, version int, reviewer string, state string) error
	AnyTasksInState(tx *sqlx.Tx, did string, DocumentNumber string, version int, states ...string) (bool, error)
	TaskExistsInState(tx *sqlx.Tx, did string, DocumentNumber string, version int, reviewer string, state string) (bool, error)
	TaskExist(tx *sqlx.Tx, did string, DocumentNumber string, version int) (bool, error)
	Delete(tx *sqlx.Tx, did string, DocumentNumber string, version int) error
}
