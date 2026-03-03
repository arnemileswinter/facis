package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	templaterepository2 "digital-contracting-service/internal/templaterepository/datatype/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/datatype/templatetype"
	"digital-contracting-service/internal/templaterepository/db"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type UpdateManageCmd struct {
	DID            string
	DocumentNumber int
	Version        int
	State          *templatestate.TemplateState
	TemplateType   *templatetype.TemplateType
	UpdatedAt      time.Time
	UpdatedBy      string
	Name           *string
	Description    *string
	TemplateData   *datatype.JSON
	IsManager      bool
}

type UpdateManager struct {
	Ctx    context.Context
	DB     *sqlx.DB
	CTRepo db.TemplateRepository
	RTRepo db.ReviewTaskRepo
	ATRepo db.ApprovalTaskRepo
}

func (h *UpdateManager) Handle(cmd UpdateManageCmd) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	oldData, err := h.CTRepo.ReadDataByID(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not read template data: %w", err)
	}

	if cmd.UpdatedAt.Before(oldData.UpdatedAt) {
		return errors.New("contract template was updated elsewhere, please reload")
	}

	if oldData.State == templatestate.Approved || oldData.State == templatestate.Registered || oldData.State == templatestate.Archived {
		return errors.New("invalid contract template state")
	}

	if cmd.State != nil {
		isValidState := *cmd.State == templatestate.Draft || *cmd.State == templatestate.Archived
		if oldData.State == templatestate.Draft && !isValidState {
			reviewTasksExist, err := h.RTRepo.TaskExist(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
			if err != nil {
				return fmt.Errorf("could not check existing review tasks: %w", err)
			}

			approvalTaskExists, err := h.ATRepo.TaskExists(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
			if err != nil {
				return fmt.Errorf("could not check existing approval tasks: %w", err)
			}

			if !reviewTasksExist || !approvalTaskExists {
				return errors.New("invalid state change")
			}
		}
	}

	newState := oldData.State
	if cmd.State != nil {
		if *cmd.State == templatestate.Draft || *cmd.State == templatestate.Archived {

			err = h.RTRepo.Delete(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
			if err != nil {
				return fmt.Errorf("could not delete review tasks: %w", err)
			}

			err = h.ATRepo.Delete(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
			if err != nil {
				return fmt.Errorf("could not delete approval tasks: %w", err)
			}

		} else if *cmd.State == templatestate.Rejected || *cmd.State == templatestate.Submitted || *cmd.State == templatestate.Reviewed {
			err = h.RTRepo.ReopenTasks(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
			if err != nil {
				return err
			}

			err = h.ATRepo.ReopenTasks(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
			if err != nil {
				return err
			}
		} else {
			return errors.New("contract invalid state")
		}

		newState = *cmd.State
	}

	newData := templaterepository2.UpdateData{
		DID:            cmd.DID,
		DocumentNumber: cmd.DocumentNumber,
		Version:        cmd.Version,
		State:          cmd.State,
		TemplateType:   cmd.TemplateType,
		Name:           cmd.Name,
		Description:    cmd.Description,
		TemplateData:   cmd.TemplateData,
	}
	err = h.CTRepo.Update(tx, newData)
	if err != nil {
		return fmt.Errorf("could not update template data: %w", err)
	}

	evt := templateevents.UpdateManageEvent{
		DID:             cmd.DID,
		DocumentNumber:  cmd.DocumentNumber,
		Version:         cmd.Version,
		OldState:        &oldData.State,
		NewState:        &newState,
		OldName:         oldData.Name,
		NewName:         cmd.Name,
		OldDescription:  oldData.Description,
		NewDescription:  cmd.Description,
		OldTemplateData: oldData.TemplateData,
		NewTemplateData: cmd.TemplateData,
		UpdatedBy:       cmd.UpdatedBy,
		OccurredAt:      time.Now(),
	}
	err = event.Create(ctx, tx, evt)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	return tx.Commit()
}
