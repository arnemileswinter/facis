package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/approvaltask"
	"digital-contracting-service/internal/templaterepository/datatype/actionflag"
	aopprovaltaskstate "digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"digital-contracting-service/internal/templaterepository/reviewtask"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type SubmitCommand struct {
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

type SubmitHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func createTasks(ctx context.Context, tx *sqlx.Tx, processData *templaterepository.ProcessData, cmd SubmitCommand) error {
	for _, reviewer := range cmd.Reviewer {
		reviewTask := reviewtask.TaskData{
			DID:            cmd.DID,
			DocumentNumber: processData.DocumentNumber,
			Version:        processData.Version,
			Reviewer:       reviewer,
			State:          reviewtaskstate.Open,
			CreatedBy:      cmd.SubmittedBy,
		}
		_, err := reviewtask.CreateTask(ctx, tx, reviewTask)
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
	_, err := approvaltask.CreateTask(ctx, tx, data)
	if err != nil {
		return fmt.Errorf("could not create approval task: %w", err)
	}

	return nil
}

func (h *SubmitHandler) Handle(cmd SubmitCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	processData, err := templaterepository.ReadProcessData(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
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

		err := createTasks(ctx, tx, processData, cmd)
		if err != nil {
			return err
		}

		nextTemplateState = templatestate.Submitted

	} else if processData.State == templatestate.Rejected {

		if processData.CreatedBy != cmd.SubmittedBy {
			return errors.New("invalid user")
		}

		err := templaterepository.ReopenTasks(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
		if err != nil {
			return err
		}

		nextTemplateState = templatestate.Submitted

	} else if processData.State == templatestate.Submitted {

		if cmd.ActionFlag != nil {
			if *cmd.ActionFlag == actionflag.Approval {

				valid, err := reviewtask.IsValidTaskUser(ctx, tx, processData.DID, processData.DocumentNumber, processData.Version, cmd.SubmittedBy)
				if err != nil {
					return err
				}

				if !valid {
					return errors.New("invalid user")
				}

				exist, err := reviewtask.HasTaskInState(ctx, tx, processData.DID, processData.DocumentNumber, processData.Version, cmd.SubmittedBy, reviewtaskstate.Open)
				if err != nil {
					return err
				}

				if exist {
					return errors.New("contract template needs to be verified before")
				}

				err = reviewtask.UpdateTask(ctx, tx, processData.DID, processData.DocumentNumber, processData.Version, cmd.SubmittedBy, reviewtaskstate.Approved)
				if err != nil {
					return fmt.Errorf("could not update approval task: %w", err)
				}

				existOpenTasks, err := reviewtask.ExistTasksInStates(ctx, tx, processData.DID, processData.DocumentNumber, processData.Version, reviewtaskstate.Open, reviewtaskstate.Verified)
				if err != nil {
					return fmt.Errorf("could not check if review task exists: %w", err)
				}

				if !existOpenTasks {
					nextTemplateState = templatestate.Reviewed
				}

			} else if *cmd.ActionFlag == actionflag.Draft {

				isValid, err := reviewtask.IsValidTaskUser(ctx, tx, processData.DID, processData.DocumentNumber, processData.Version, cmd.SubmittedBy)
				if err != nil {
					return err
				}

				if !isValid {
					return errors.New("invalid user")
				}

				err = templaterepository.ReopenTasks(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
				if err != nil {
					return err
				}

				nextTemplateState = templatestate.Rejected
			}
		} else {
			return errors.New("action flags is missing")
		}

	} else if processData.State == templatestate.Reviewed {

		isValid, err := approvaltask.IsValidTaskUser(ctx, tx, processData.DID, processData.DocumentNumber, processData.Version, cmd.SubmittedBy)
		if err != nil {
			return err
		}

		if !isValid {
			return errors.New("invalid user")
		}

		err = templaterepository.ReopenTasks(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
		if err != nil {
			return err
		}

		nextTemplateState = templatestate.Submitted

	} else {
		return errors.New("current contract template state is invalid")
	}

	if len(nextTemplateState) > 0 && processData.State != nextTemplateState {
		err = templaterepository.UpdateState(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, nextTemplateState)
		if err != nil {
			return fmt.Errorf("could not update contract template state: %w", err)
		}

		evt := templateevents.SubmitContractTemplateEvent{
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
