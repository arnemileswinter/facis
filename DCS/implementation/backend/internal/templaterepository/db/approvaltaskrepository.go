package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type ApprovalTaskData struct {
	ID             string    `db:"id"`
	DID            string    `db:"did"`
	DocumentNumber int       `db:"document_number"`
	Version        int       `db:"version"`
	State          string    `db:"state"`
	Approver       string    `db:"approver"`
	CreatedBy      string    `db:"created_by"`
	CreatedAt      time.Time `db:"created_at"`
}

type ApprovalTaskRepo interface {
	Create(tx *sqlx.Tx, data ApprovalTaskData) (*time.Time, error)
	ReopenTasks(tx *sqlx.Tx, did string, documentNumber int, version int) error
	ReadAll(dtx *sqlx.Tx, id string) ([]ApprovalTaskData, error)
	ReadAllByApprover(tx *sqlx.Tx, approver string) ([]ApprovalTaskData, error)
	Update(tx *sqlx.Tx, did string, documentNumber int, version int, approver string, state string) error
	IsValidApprover(tx *sqlx.Tx, did string, documentNumber int, version int, approver string) (bool, error)
	TaskExistsInState(tx *sqlx.Tx, did string, documentNumber int, version int, approver string, state string) (bool, error)
	TaskExists(tx *sqlx.Tx, did string, documentNumber int, version int) (bool, error)
	Delete(tx *sqlx.Tx, did string, documentNumber int, version int) error
}
