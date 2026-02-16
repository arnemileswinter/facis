package query

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	aopprovaltaskstate "digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"time"

	"github.com/jmoiron/sqlx"
)

type GetAllContractTemplateApprovalTasksForDID struct {
	DID            string
	DocumentNumber int
	Version        int
	RetrievedBy    string
}

type GetAllContractTemplateApprovalTasksForDIDResult struct {
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

type GetAllContractTemplateApprovalTasksForDIDHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *GetAllContractTemplateApprovalTasksForDIDHandler) Handle(query GetAllContractTemplateApprovalTasksForDID) ([]GetAllContractTemplateApprovalTasksForDIDResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, base.GetTransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	reviewTasks, err := templaterepository.ReadAllApprovalTasks(ctx, tx, query.DID)
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

	result := make([]GetAllContractTemplateApprovalTasksForDIDResult, len(reviewTasks))
	for i, data := range reviewTasks {
		result[i] = GetAllContractTemplateApprovalTasksForDIDResult{
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
