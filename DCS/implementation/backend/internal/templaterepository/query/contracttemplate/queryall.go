package contracttemplate

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/contracttemplatestate"
	"digital-contracting-service/internal/templaterepository/datatype/contracttemplatetype"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/db"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type GetAllMetadataQry struct {
	RetrievedBy string
}

type MetadataItem struct {
	DID            string
	DocumentNumber int
	Version        int
	State          contracttemplatestate.ContractTemplateState
	TemplateType   contracttemplatetype.ContractTemplateType
	Name           string
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	MetaData       datatype.JSON
}

type ReviewTaskItem struct {
	DID            string
	DocumentNumber int
	Version        int
	State          reviewtaskstate.ReviewTaskState
	Reviewer       string
	CreatedAt      time.Time
}

type ApprovalTaskItem struct {
	DID            string
	DocumentNumber int
	Version        int
	State          approvaltaskstate.ApprovalTaskState
	Approver       string
	CreatedAt      time.Time
}

type GetAllMetadataResult struct {
	ContractTemplates []MetadataItem
	ReviewerTasks     []ReviewTaskItem
	ApprovalTasks     []ApprovalTaskItem
}

type GetAllMetadataHandler struct {
	Ctx    context.Context
	DB     *sqlx.DB
	CTRepo db.ContractTemplateRepo
	RTRepo db.ReviewTaskRepo
	ATRepo db.ApprovalTaskRepo
}

func (h *GetAllMetadataHandler) Handle(query GetAllMetadataQry) (*GetAllMetadataResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create transaction: %w", err)
	}
	defer tx.Rollback()

	contractTemplates, err := h.CTRepo.ReadAllMetaData(tx)
	if err != nil {
		return nil, fmt.Errorf("could not read all contract templates: %w", err)
	}

	evt := templateevents.RetrieveAllEvent{
		RetrievedBy: query.RetrievedBy,
		OccurredAt:  time.Now(),
	}
	err = event.Create(h.Ctx, tx, evt)
	if err != nil {
		return nil, fmt.Errorf("could not create event: %w", err)
	}

	reviewerTasks, err := h.RTRepo.ReadAllByReviewer(tx, query.RetrievedBy)
	if err != nil {
		return nil, fmt.Errorf("could not read all review tasks: %w", err)
	}

	approvalTasks, err := h.ATRepo.ReadAllByApprover(tx, query.RetrievedBy)
	if err != nil {
		return nil, fmt.Errorf("could not read all review tasks: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	var contractTemplatesItems []MetadataItem
	for _, data := range contractTemplates {

		state, err := contracttemplatestate.NewContractTemplateState(data.State)
		if err != nil {
			return nil, fmt.Errorf("could not create contract template state: %w", err)
		}

		templateType, err := contracttemplatetype.NewContractTemplateType(data.TemplateType)
		if err != nil {
			return nil, fmt.Errorf("could not create contract template type: %w", err)
		}

		contractTemplatesItems = append(contractTemplatesItems, MetadataItem{
			DID:            data.DID,
			DocumentNumber: data.DocumentNumber,
			Version:        data.Version,
			State:          state,
			TemplateType:   templateType,
			Name:           *data.Name,
			Description:    *data.Description,
			CreatedAt:      data.CreatedAt,
			UpdatedAt:      data.UpdatedAt,
		})
	}

	var reviewTaskItems []ReviewTaskItem
	for _, data := range reviewerTasks {

		state, err := reviewtaskstate.NewReviewTaskState(data.State)
		if err != nil {
			return nil, fmt.Errorf("could not create contract template state: %w", err)
		}

		reviewTaskItems = append(reviewTaskItems, ReviewTaskItem{
			DID:            data.DID,
			DocumentNumber: data.DocumentNumber,
			Version:        data.Version,
			State:          state,
			Reviewer:       data.Reviewer,
			CreatedAt:      data.CreatedAt,
		})
	}

	var approvalTasksItems []ApprovalTaskItem
	for _, data := range approvalTasks {

		state, err := approvaltaskstate.NewApprovalTaskState(data.State)
		if err != nil {
			return nil, fmt.Errorf("could not create contract template state: %w", err)
		}

		approvalTasksItems = append(approvalTasksItems, ApprovalTaskItem{
			DID:            data.DID,
			DocumentNumber: data.DocumentNumber,
			Version:        data.Version,
			State:          state,
			Approver:       data.Approver,
			CreatedAt:      data.CreatedAt,
		})
	}

	return &GetAllMetadataResult{
		ContractTemplates: contractTemplatesItems,
		ReviewerTasks:     reviewTaskItems,
		ApprovalTasks:     approvalTasksItems,
	}, nil
}
