package contracttemplate

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/approvaltask"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"digital-contracting-service/internal/templaterepository/reviewtask"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type GetAllContractTemplatesMetaData struct {
	RetrievedBy string
}

type ContractTemplatesMetaDataItem struct {
	DID            string
	DocumentNumber int
	Version        int
	State          templatestate.TemplateState
	Name           string
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	MetaData       datatype.JSON
}

type ContractTemplatesReviewTaskItem struct {
	DID            string
	DocumentNumber int
	Version        int
	State          reviewtaskstate.ReviewTaskState
	Reviewer       string
	CreatedAt      time.Time
}

type ContractTemplatesApprovalTaskItem struct {
	DID            string
	DocumentNumber int
	Version        int
	State          approvaltaskstate.ApprovalTaskState
	Approver       string
	CreatedAt      time.Time
}

type GetAllContractTemplatesMetaDataResult struct {
	ContractTemplates []ContractTemplatesMetaDataItem
	ReviewerTasks     []ContractTemplatesReviewTaskItem
	ApprovalTasks     []ContractTemplatesApprovalTaskItem
}

type GetAllContractTemplateMetaDataHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *GetAllContractTemplateMetaDataHandler) Handle(query GetAllContractTemplatesMetaData) (*GetAllContractTemplatesMetaDataResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create transaction: %w", err)
	}
	defer tx.Rollback()

	contractTemplates, err := templaterepository.ReadAllMetaData(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("could not read all contract templates: %w", err)
	}

	evt := templateevents.RetrieveAllContractTemplatesEvent{
		RetrievedBy: query.RetrievedBy,
		OccurredAt:  time.Now(),
	}
	err = event.Create(h.Ctx, tx, evt)
	if err != nil {
		return nil, fmt.Errorf("could not create event: %w", err)
	}

	reviewerTasks, err := reviewtask.ReadAllByReviewer(ctx, tx, query.RetrievedBy)
	if err != nil {
		return nil, fmt.Errorf("could not read all review tasks: %w", err)
	}

	approvalTasks, err := approvaltask.ReadAllByApprover(ctx, tx, query.RetrievedBy)
	if err != nil {
		return nil, fmt.Errorf("could not read all review tasks: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	var contractTemplatesItems []ContractTemplatesMetaDataItem
	for _, data := range contractTemplates {
		contractTemplatesItems = append(contractTemplatesItems, ContractTemplatesMetaDataItem{
			DID:            data.DID,
			DocumentNumber: data.DocumentNumber,
			Version:        data.Version,
			State:          data.State,
			Name:           *data.Name,
			Description:    *data.Description,
			CreatedAt:      data.CreatedAt,
			UpdatedAt:      data.UpdatedAt,
		})
	}

	var reviewTaskItems []ContractTemplatesReviewTaskItem
	for _, data := range reviewerTasks {
		reviewTaskItems = append(reviewTaskItems, ContractTemplatesReviewTaskItem{
			DID:            data.DID,
			DocumentNumber: data.DocumentNumber,
			Version:        data.Version,
			State:          data.State,
			Reviewer:       data.Reviewer,
			CreatedAt:      data.CreatedAt,
		})
	}

	var approvalTasksItems []ContractTemplatesApprovalTaskItem
	for _, data := range approvalTasks {
		approvalTasksItems = append(approvalTasksItems, ContractTemplatesApprovalTaskItem{
			DID:            data.DID,
			DocumentNumber: data.DocumentNumber,
			Version:        data.Version,
			State:          data.State,
			Approver:       data.Approver,
			CreatedAt:      data.CreatedAt,
		})
	}

	return &GetAllContractTemplatesMetaDataResult{
		ContractTemplates: contractTemplatesItems,
		ReviewerTasks:     reviewTaskItems,
		ApprovalTasks:     approvalTasksItems,
	}, nil
}
