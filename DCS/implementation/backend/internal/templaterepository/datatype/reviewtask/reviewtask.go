package reviewtask

import (
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"time"
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
