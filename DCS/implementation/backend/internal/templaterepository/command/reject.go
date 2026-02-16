package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type RejectTemplateContractCommand struct {
	DID            string
	DocumentNumber int
	Version        int
	UpdatedAt      time.Time
	RejectedBy     string
	Reason         string
}

type RejectTemplateContractHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *RejectTemplateContractHandler) Handle(cmd RejectTemplateContractCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	processData, err := templaterepository.ReadContractTemplateProcessData(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not read process data: %w", err)
	}

	if cmd.UpdatedAt.Before(processData.UpdatedAt) {
		return errors.New("contract template was updated elsewhere, please reload")
	}

	if processData.State != templatestate.Reviewed {
		return errors.New("invalid contract template state")
	}

	err = templaterepository.UpdateContractTemplateState(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, templatestate.Draft)
	if err != nil {
		return fmt.Errorf("could not update current template state: %w", err)
	}

	evt := templateevents.ContractTemplateRejectedEvent{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		RejectedBy:     cmd.RejectedBy,
		Reason:         cmd.Reason,
		OccurredAt:     time.Now(),
	}
	err = event.Create(ctx, tx, evt)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	err = templaterepository.DeleteReviewTask(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not delete review task: %w", err)
	}

	deleteReviewTaskEvt := templateevents.ContractTemplateDeleteReviewTaskEvent{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		DeletedBy:      cmd.RejectedBy,
		OccurredAt:     time.Now(),
	}
	err = event.Create(ctx, tx, deleteReviewTaskEvt)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	err = templaterepository.DeleteApprovalTask(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not delete approval task: %w", err)
	}

	deleteApprovalTaskEvt := templateevents.ContractTemplateDeleteApprovalTaskEvent{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		DeletedBy:      cmd.RejectedBy,
		OccurredAt:     time.Now(),
	}
	err = event.Create(ctx, tx, deleteApprovalTaskEvt)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	return tx.Commit()
}
