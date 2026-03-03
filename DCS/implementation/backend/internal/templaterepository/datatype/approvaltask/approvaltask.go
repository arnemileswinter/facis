package approvaltask

import (
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"time"
)

type TaskData struct {
	ID             string                              `db:"id"`
	DID            string                              `db:"did"`
	DocumentNumber int                                 `db:"document_number"`
	Version        int                                 `db:"version"`
	State          approvaltaskstate.ApprovalTaskState `db:"state"`
	Approver       string                              `db:"approver"`
	CreatedBy      string                              `db:"created_by"`
	CreatedAt      time.Time                           `db:"created_at"`
}
