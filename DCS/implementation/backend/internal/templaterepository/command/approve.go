package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/approvaltask"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"digital-contracting-service/internal/templaterepository/reviewtask"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type ApproveCommand struct {
	DID            string
	DocumentNumber int
	Version        int
	UpdatedAt      time.Time
	ApprovedBy     string
	DecisionNotes  []string
}

type ApproveHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *ApproveHandler) Handle(cmd ApproveCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	processData, err := templaterepository.ReadProcessData(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not read process data: %w", err)
	}

	if cmd.UpdatedAt.Before(processData.UpdatedAt) {
		return errors.New("contract template was updated elsewhere, please reload")
	}

	if processData.State != templatestate.Reviewed {
		return errors.New("invalid contract template state")
	}

	valid, err := approvaltask.IsValidTaskUser(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, cmd.ApprovedBy)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid user")
	}

	exist, err := approvaltask.HasTaskInState(ctx, tx, processData.DID, processData.DocumentNumber, processData.Version, cmd.ApprovedBy, approvaltaskstate.Open)
	if err != nil {
		return err
	}

	if exist {
		return errors.New("contract template needs to be verified before")
	}

	err = templaterepository.UpdateState(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, templatestate.Approved)
	if err != nil {
		return fmt.Errorf("could not update current template state: %w", err)
	}

	evt := templateevents.ApproveContractTemplateEvent{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		ApprovedBy:     cmd.ApprovedBy,
		DecisionNotes:  cmd.DecisionNotes,
		OccurredAt:     time.Now(),
	}
	err = event.Create(ctx, tx, evt)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	err = reviewtask.DeleteTask(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not delete review task: %w", err)
	}

	deleteReviewTaskEvt := templateevents.DeleteContractTemplateReviewTaskEvent{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		DeletedBy:      cmd.ApprovedBy,
		OccurredAt:     time.Now(),
	}
	err = event.Create(ctx, tx, deleteReviewTaskEvt)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	err = approvaltask.DeleteTask(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not delete approval task: %w", err)
	}

	deleteApprovalTaskEvt := templateevents.DeleteContractTemplateApprovalTaskEvent{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		DeletedBy:      cmd.ApprovedBy,
		OccurredAt:     time.Now(),
	}
	err = event.Create(ctx, tx, deleteApprovalTaskEvt)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	return tx.Commit()
}
