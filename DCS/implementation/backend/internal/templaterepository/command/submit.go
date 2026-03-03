package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository/datatype/actionflag"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltask"
	aopprovaltaskstate "digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtask"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	templaterepository2 "digital-contracting-service/internal/templaterepository/datatype/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/db"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type SubmitCmd struct {
	DID            string
	DocumentNumber int
	Version        int
	UpdatedAt      time.Time
	SubmittedBy    string
	ActionFlag     *actionflag.ActionFlag
	Comments       []string
	Reviewer       []string
	Approver       *string
}

type Submitter struct {
	Ctx    context.Context
	DB     *sqlx.DB
	CTRepo db.TemplateRepository
	RTRepo db.ReviewTaskRepo
	ATRepo db.ApprovalTaskRepo
}

func createTasks(tx *sqlx.Tx, rtRepo db.ReviewTaskRepo, atRepo db.ApprovalTaskRepo, processData *templaterepository2.ProcessData, cmd SubmitCmd) error {
	for _, reviewer := range cmd.Reviewer {
		reviewTask := reviewtask.TaskData{
			DID:            cmd.DID,
			DocumentNumber: processData.DocumentNumber,
			Version:        processData.Version,
			Reviewer:       reviewer,
			State:          reviewtaskstate.Open,
			CreatedBy:      cmd.SubmittedBy,
		}
		_, err := rtRepo.Create(tx, reviewTask)
		if err != nil {
			return fmt.Errorf("could not create review tasks: %w", err)
		}
	}

	data := approvaltask.TaskData{
		DID:            cmd.DID,
		DocumentNumber: processData.DocumentNumber,
		Version:        processData.Version,
		CreatedBy:      cmd.SubmittedBy,
		Approver:       *cmd.Approver,
		State:          aopprovaltaskstate.Open,
	}
	_, err := atRepo.Create(tx, data)
	if err != nil {
		return fmt.Errorf("could not create approval task: %w", err)
	}

	return nil
}

func (h *Submitter) Handle(cmd SubmitCmd) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	processData, err := h.CTRepo.ReadProcessData(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return fmt.Errorf("could not process core data: %w", err)
	}

	if cmd.UpdatedAt.Before(processData.UpdatedAt) {
		return errors.New("contract template was updated elsewhere, please reload")
	}

	var nextTemplateState templatestate.TemplateState
	if processData.State == templatestate.Draft {

		if cmd.SubmittedBy != processData.CreatedBy {
			return errors.New("invalid user")
		}

		if len(cmd.Reviewer) == 0 {
			return errors.New("no reviewer provided")
		}

		if cmd.Approver == nil || len(*cmd.Approver) == 0 {
			return errors.New("no approver provided")
		}

		err := createTasks(tx, h.RTRepo, h.ATRepo, processData, cmd)
		if err != nil {
			return err
		}

		nextTemplateState = templatestate.Submitted

	} else if processData.State == templatestate.Rejected {

		if processData.CreatedBy != cmd.SubmittedBy {
			return errors.New("invalid user")
		}

		err := h.RTRepo.ReopenTasks(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
		if err != nil {
			return errors.New("could not reopen review tasks")
		}

		err = h.ATRepo.ReopenTasks(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
		if err != nil {
			return errors.New("could not reopen approval tasks")
		}

		nextTemplateState = templatestate.Submitted

	} else if processData.State == templatestate.Submitted {

		if cmd.ActionFlag != nil {
			if *cmd.ActionFlag == actionflag.Approval {

				valid, err := h.RTRepo.IsValidReviewer(tx, processData.DID, processData.DocumentNumber, processData.Version, cmd.SubmittedBy)
				if err != nil {
					return err
				}

				if !valid {
					return errors.New("invalid user")
				}

				exist, err := h.RTRepo.TaskExistsInState(tx, processData.DID, processData.DocumentNumber, processData.Version, cmd.SubmittedBy, reviewtaskstate.Open)
				if err != nil {
					return err
				}

				if exist {
					return errors.New("contract template needs to be verified before")
				}

				err = h.RTRepo.Update(tx, processData.DID, processData.DocumentNumber, processData.Version, cmd.SubmittedBy, reviewtaskstate.Approved)
				if err != nil {
					return fmt.Errorf("could not update approval task: %w", err)
				}

				existOpenTasks, err := h.RTRepo.AnyTasksInState(tx, processData.DID, processData.DocumentNumber, processData.Version, reviewtaskstate.Open, reviewtaskstate.Verified)
				if err != nil {
					return fmt.Errorf("could not check if review task exists: %w", err)
				}

				if !existOpenTasks {
					nextTemplateState = templatestate.Reviewed
				}

			} else if *cmd.ActionFlag == actionflag.Draft {

				isValid, err := h.RTRepo.IsValidReviewer(tx, processData.DID, processData.DocumentNumber, processData.Version, cmd.SubmittedBy)
				if err != nil {
					return err
				}

				if !isValid {
					return errors.New("invalid user")
				}

				err = h.RTRepo.ReopenTasks(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
				if err != nil {
					return err
				}

				err = h.ATRepo.ReopenTasks(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
				if err != nil {
					return err
				}

				nextTemplateState = templatestate.Rejected
			}
		} else {
			return errors.New("action flags is missing")
		}

	} else if processData.State == templatestate.Reviewed {

		isValid, err := h.ATRepo.IsValidApprover(tx, processData.DID, processData.DocumentNumber, processData.Version, cmd.SubmittedBy)
		if err != nil {
			return err
		}

		if !isValid {
			return errors.New("invalid user")
		}

		err = h.RTRepo.ReopenTasks(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
		if err != nil {
			return err
		}

		err = h.ATRepo.ReopenTasks(tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
		if err != nil {
			return err
		}
		nextTemplateState = templatestate.Submitted

	} else {
		return errors.New("current contract template state is invalid")
	}

	if len(nextTemplateState) > 0 && processData.State != nextTemplateState {
		err = h.CTRepo.UpdateState(tx, cmd.DID, cmd.DocumentNumber, cmd.Version, nextTemplateState)
		if err != nil {
			return fmt.Errorf("could not update contract template state: %w", err)
		}

		evt := templateevents.SubmitEvent{
			DID:            cmd.DID,
			DocumentNumber: cmd.DocumentNumber,
			Version:        cmd.DocumentNumber,
			SubmittedBy:    cmd.SubmittedBy,
			PreviousState:  processData.State,
			NewState:       nextTemplateState,
			ActionFlag:     cmd.ActionFlag,
			Comments:       cmd.Comments,
			OccurredAt:     time.Now(),
		}
		err = event.Create(ctx, tx, evt)
		if err != nil {
			return fmt.Errorf("could not create event: %w", err)
		}
	}

	return tx.Commit()
}
