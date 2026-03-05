package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository/datatype/contracttemplatestate"
	"digital-contracting-service/internal/templaterepository/db"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type ArchiveCmd struct {
	DID            string
	DocumentNumber int
	Version        int
	UpdatedAt      time.Time
	ArchivedBy     string
}

type Archiver struct {
	Ctx    context.Context
	DB     *sqlx.DB
	CTRepo db.ContractTemplateRepo
	RTRepo db.ReviewTaskRepo
	ATRepo db.ApprovalTaskRepo
}

func (h *Archiver) Handle(cmd ArchiveCmd) error {

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

	if processData.State == contracttemplatestate.Deprecated.String() || processData.State == contracttemplatestate.Deleted.String() {
		return errors.New("invalid contract template state")
	}

	if processData.State == contracttemplatestate.Registered.String() {

		err = h.CTRepo.UpdateState(tx, cmd.DID, cmd.DocumentNumber, cmd.Version, contracttemplatestate.Deprecated.String())
		if err != nil {
			return fmt.Errorf("could not update state: %w", err)
		}

	} else {

		err = h.CTRepo.UpdateState(tx, cmd.DID, cmd.DocumentNumber, cmd.Version, contracttemplatestate.Deleted.String())
		if err != nil {
			return fmt.Errorf("could not update state: %w", err)
		}
	}

	evt := templateevents.ArchiveEvent{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		ArchivedBy:     cmd.ArchivedBy,
		OccurredAt:     time.Now(),
	}
	err = event.Create(ctx, tx, evt)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	err = h.RTRepo.Delete(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not delete review tasks: %w", err)
	}

	err = h.ATRepo.Delete(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not delete approval tasks: %w", err)
	}

	return tx.Commit()
}
