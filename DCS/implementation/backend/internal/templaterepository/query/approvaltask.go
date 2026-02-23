package query

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository/approvaltask"
	aopprovaltaskstate "digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type GetAllApprovalTasksForDID struct {
	DID            string
	DocumentNumber int
	Version        int
	RetrievedBy    string
}

type GetAllApprovalTasksForDIDResult struct {
	ID             int
	DID            string
	DocumentNumber int
	Version        int
	State          aopprovaltaskstate.ApprovalTaskState
	Approver       string
	CreatedBy      string
	CreatedAt      time.Time
	CancelledAt    *time.Time
}

type GetAllApprovalTasksForDIDHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *GetAllApprovalTasksForDIDHandler) Handle(query GetAllApprovalTasksForDID) ([]GetAllApprovalTasksForDIDResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	reviewTasks, err := approvaltask.ReadAll(ctx, tx, query.DID)
	if err != nil {
		return nil, fmt.Errorf("could not read all review tasks: %w", err)
	}

	evt := templateevents.RetrieveAllContractTemplateReviewTasksEvent{
		DID:            query.DID,
		DocumentNumber: query.DocumentNumber,
		Version:        query.Version,
		RetrievedBy:    query.RetrievedBy,
		OccurredAt:     time.Now(),
	}
	err = event.Create(h.Ctx, tx, evt)
	if err != nil {
		return nil, fmt.Errorf("could not create event: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	result := make([]GetAllApprovalTasksForDIDResult, len(reviewTasks))
	for i, data := range reviewTasks {
		result[i] = GetAllApprovalTasksForDIDResult{
			DID:            data.DID,
			DocumentNumber: data.DocumentNumber,
			Version:        data.Version,
			State:          data.State,
			Approver:       data.Approver,
			CreatedBy:      data.CreatedBy,
			CreatedAt:      data.CreatedAt,
		}
	}

	return result, nil
}
