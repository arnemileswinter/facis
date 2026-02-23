package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type VerifyCommand struct {
	DID            string
	DocumentNumber int
	Version        int
	UpdatedAt      time.Time
	VerifiedBy     string
}

type VerifyHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *VerifyHandler) Handle(cmd VerifyCommand) error {

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

	fmt.Println("process data:", processData)

	evt := templateevents.VerifyContractTemplateEvent{
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
