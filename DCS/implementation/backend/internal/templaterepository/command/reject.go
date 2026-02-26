package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/approvaltask"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type RejectCommand struct {
	DID            string
	DocumentNumber int
	Version        int
	UpdatedAt      time.Time
	RejectedBy     string
	Reason         string
}

type RejectHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *RejectHandler) Handle(cmd RejectCommand) error {

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

	exist, err := approvaltask.IsValidTaskUser(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, cmd.RejectedBy)
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("invalid user")
	}

	err = templaterepository.UpdateState(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, templatestate.Draft)
	if err != nil {
		return fmt.Errorf("could not update current template state: %w", err)
	}

	evt := templateevents.RejectEvent{
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

	err = templaterepository.CleanupTasks(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not cleanup tasks: %w", err)
	}

	return tx.Commit()
}
