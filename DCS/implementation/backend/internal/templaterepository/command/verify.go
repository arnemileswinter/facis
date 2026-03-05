package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype/componenttype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/db"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type VerifyCmd struct {
	DID            string
	DocumentNumber string
	Version        int
	UpdatedAt      time.Time
	VerifiedBy     string
}

type Verifier struct {
	Ctx    context.Context
	DB     *sqlx.DB
	CTRepo db.ContractTemplateRepo
	RTRepo db.ReviewTaskRepo
	ATRepo db.ApprovalTaskRepo
}

func (h *Verifier) Handle(cmd VerifyCmd) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	processData, err := h.CTRepo.ReadProcessData(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not read process data: %w", err)
	}

	if cmd.UpdatedAt.Before(processData.UpdatedAt) {
		return errors.New("contract template was updated elsewhere, please reload")
	}

	hasTask, err := h.RTRepo.TaskExistsInState(tx, cmd.DID, cmd.DocumentNumber, cmd.Version, cmd.VerifiedBy, reviewtaskstate.Open.String())
	if err != nil {
		return err
	}

	if hasTask {
		err := h.RTRepo.Update(tx, cmd.DID, cmd.DocumentNumber, cmd.Version, cmd.VerifiedBy, reviewtaskstate.Verified.String())
		if err != nil {
			return err
		}
	}

	hasTask, err = h.ATRepo.TaskExistsInState(tx, cmd.DID, cmd.DocumentNumber, cmd.Version, cmd.VerifiedBy, reviewtaskstate.Open.String())
	if err != nil {
		return err
	}

	if hasTask {
		err := h.ATRepo.Update(tx, cmd.DID, cmd.DocumentNumber, cmd.Version, cmd.VerifiedBy, reviewtaskstate.Verified.String())
		if err != nil {
			return err
		}
	}

	evt := templateevents.VerifyEvent{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		VerifiedBy:     cmd.VerifiedBy,
		OccurredAt:     time.Now(),
	}
	err = event.Create(ctx, tx, evt, componenttype.ContractTemplateRepo)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	return tx.Commit()
}
