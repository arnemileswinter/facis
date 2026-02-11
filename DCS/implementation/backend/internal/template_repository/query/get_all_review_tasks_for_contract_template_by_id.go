package query

import (
	"context"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/template_repository"
	"digital-contracting-service/internal/template_repository/datatype/review_task_state"
	templateevents "digital-contracting-service/internal/template_repository/event"
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
	State          review_task_state.ReviewTaskState
	Reviewer       string
	CreatedBy      string
	CreatedAt      time.Time
	CancelledAt    *time.Time
}

type GetAllContractTemplateReviewTasksForDIDHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *GetAllContractTemplateReviewTasksForDIDHandler) Handle(query GetAllContractTemplateReviewTasksForDID) ([]GetAllContractTemplateReviewTasksForDIDResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	reviewTasks, err := template_repository.ReadAllReviewTasks(ctx, tx, query.DID)
	if err != nil {
		return nil, err
	}

	evt := templateevents.ContractTemplateRetrieveAllReviewTasksEvent{
		DID:            query.DID,
		DocumentNumber: query.DocumentNumber,
		Version:        query.Version,
		RetrievedBy:    query.RetrievedBy,
		OccurredAt:     time.Now(),
	}
	err = event.CreateNewEvent(h.Ctx, tx, evt)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
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
			CancelledAt:    data.CancelledAt,
		}
	}

	return result, nil
}
