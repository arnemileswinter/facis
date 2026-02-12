package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/template_repository"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	templateevents "digital-contracting-service/internal/template_repository/event"
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

	currentTemplateState, err := template_repository.ReadContractTemplateState(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return err
	}

	if *currentTemplateState != template_state.Reviewed {
		return errors.New("invalid contract template state")
	}

	err = template_repository.UpdateContractTemplateState(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, template_state.Approved)
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

	return tx.Commit()
}
