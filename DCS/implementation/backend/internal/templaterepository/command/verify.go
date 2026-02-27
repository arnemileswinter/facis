package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/approvaltask"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"digital-contracting-service/internal/templaterepository/reviewtask"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type VerifyCmd struct {
	DID            string
	DocumentNumber int
	Version        int
	UpdatedAt      time.Time
	VerifiedBy     string
}

type Verifier struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *Verifier) Handle(cmd VerifyCmd) error {

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

	hasTask, err := reviewtask.TaskExistsInState(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, cmd.VerifiedBy, reviewtaskstate.Open)
	if err != nil {
		return err
	}

	if hasTask {
		err := reviewtask.Update(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, cmd.VerifiedBy, reviewtaskstate.Verified)
		if err != nil {
			return err
		}
	}

	hasTask, err = approvaltask.TaskExistsInState(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, cmd.VerifiedBy, approvaltaskstate.Open)
	if err != nil {
		return err
	}

	if hasTask {
		err := approvaltask.Update(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, cmd.VerifiedBy, approvaltaskstate.Verified)
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
	err = event.Create(ctx, tx, evt)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	return tx.Commit()
}
