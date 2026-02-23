package query

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"digital-contracting-service/internal/templaterepository/reviewtask"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type GetAllContractTemplateReviewTasksForDID struct {
	DID            string
	DocumentNumber int
	Version        int
	RetrievedBy    string
}

type GetAllContractTemplateReviewTasksForDIDResult struct {
	ID             int
	DID            string
	DocumentNumber int
	Version        int
	State          reviewtaskstate.ReviewTaskState
	Reviewer       string
	CreatedBy      string
	CreatedAt      time.Time
}

type GetAllContractTemplateReviewTasksForDIDHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *GetAllContractTemplateReviewTasksForDIDHandler) Handle(query GetAllContractTemplateReviewTasksForDID) ([]GetAllContractTemplateReviewTasksForDIDResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	reviewTasks, err := reviewtask.ReadAll(ctx, tx, query.DID)
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

	result := make([]GetAllContractTemplateReviewTasksForDIDResult, len(reviewTasks))
	for i, data := range reviewTasks {
		result[i] = GetAllContractTemplateReviewTasksForDIDResult{
			DID:            data.DID,
			DocumentNumber: data.DocumentNumber,
			Version:        data.Version,
			State:          data.State,
			Reviewer:       data.Reviewer,
			CreatedBy:      data.CreatedBy,
			CreatedAt:      data.CreatedAt,
		}
	}

	return result, nil
}
