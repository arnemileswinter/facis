package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type ApproveTemplateContractCommand struct {
	DID            string
	DocumentNumber int
	Version        int
	ApprovedBy     string
	DecisionNotes  []string
}

type ApproveTemplateContractHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *ApproveTemplateContractHandler) Handle(cmd ApproveTemplateContractCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.GetTransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	currentTemplateState, err := templaterepository.ReadContractTemplateState(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return err
	}

	if *currentTemplateState != templatestate.Reviewed {
		return errors.New("invalid contract template state")
	}

	err = templaterepository.UpdateContractTemplateState(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, templatestate.Approved)
	if err != nil {
		return err
	}

	evt := templateevents.ContractTemplateApprovedEvent{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		ApprovedBy:     cmd.ApprovedBy,
		DecisionNotes:  cmd.DecisionNotes,
		OccurredAt:     time.Now(),
	}
	err = event.CreateNewEvent(ctx, tx, evt)
	if err != nil {
		return err
	}

	err = templaterepository.DeleteReviewTask(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return err
	}

	deleteReviewTaskEvt := templateevents.ContractTemplateDeleteReviewTaskEvent{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		DeletedBy:      cmd.ApprovedBy,
		OccurredAt:     time.Now(),
	}
	err = event.CreateNewEvent(ctx, tx, deleteReviewTaskEvt)
	if err != nil {
		return err
	}

	err = templaterepository.DeleteApprovalTask(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return err
	}

	deleteApprovalTaskEvt := templateevents.ContractTemplateDeleteApprovalTaskEvent{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		DeletedBy:      cmd.ApprovedBy,
		OccurredAt:     time.Now(),
	}
	err = event.CreateNewEvent(ctx, tx, deleteApprovalTaskEvt)
	if err != nil {
		return err
	}

	return tx.Commit()
}
