package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/actionflag"
	aopprovaltaskstate "digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type SubmitContractTemplateCommand struct {
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

type SubmitContractTemplateHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func reopenReviewTasks(ctx context.Context, tx *sqlx.Tx, submittedBy string, processData *templaterepository.ContractTemplateProcessData) error {

	err := templaterepository.ReopenReviewTasks(ctx, tx, processData.DID, processData.DocumentNumber, processData.Version)
	if err != nil {
		return fmt.Errorf("could not reopen review tasks: %w", err)
	}

	reopenReviewTaskEvent := templateevents.ContractTemplateReopenReviewTaskEvent{
		DID:            processData.DID,
		DocumentNumber: processData.DocumentNumber,
		Version:        processData.Version,
		CreatedBy:      submittedBy,
		OccurredAt:     time.Now(),
	}
	err = event.Create(ctx, tx, reopenReviewTaskEvent)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	err = templaterepository.ReopenApprovalTask(ctx, tx, processData.DID, processData.DocumentNumber, processData.Version)
	if err != nil {
		return fmt.Errorf("could not reopen approval tasks: %w", err)
	}

	reopenApprovalTaskEvent := templateevents.ContractTemplateReopenApprovalTaskEvent{
		DID:            processData.DID,
		DocumentNumber: processData.DocumentNumber,
		Version:        processData.Version,
		CreatedBy:      submittedBy,
		OccurredAt:     time.Now(),
	}
	err = event.Create(ctx, tx, reopenApprovalTaskEvent)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	return nil
}

func (h *SubmitContractTemplateHandler) Handle(cmd SubmitContractTemplateCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	processData, err := templaterepository.ReadContractTemplateProcessData(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
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

		if cmd.Reviewer == nil || len(cmd.Reviewer) == 0 {
			return errors.New("no reviewer provided")
		}

		if cmd.Approver == nil || len(*cmd.Approver) == 0 {
			return errors.New("no approver provided")
		}

		for _, reviewer := range cmd.Reviewer {
			reviewTask := templaterepository.ReviewTaskData{
				DID:            cmd.DID,
				DocumentNumber: processData.DocumentNumber,
				Version:        processData.Version,
				Reviewer:       reviewer,
				State:          reviewtaskstate.Open,
				CreatedBy:      cmd.SubmittedBy,
			}
			createdAt, err := templaterepository.CreateReviewTasks(ctx, tx, reviewTask)
			if err != nil {
				return fmt.Errorf("could not create review tasks: %w", err)
			}

			createReviewTaskEvent := templateevents.ContractTemplateCreateReviewTaskEvent{
				DID:            processData.DID,
				DocumentNumber: processData.DocumentNumber,
				Version:        processData.Version,
				CreatedBy:      cmd.SubmittedBy,
				Reviewer:       reviewer,
				OccurredAt:     *createdAt,
			}
			err = event.Create(ctx, tx, createReviewTaskEvent)
			if err != nil {
				return fmt.Errorf("could not create event: %w", err)
			}
		}

		data := templaterepository.ApprovalTaskData{
			DID:            cmd.DID,
			DocumentNumber: processData.DocumentNumber,
			Version:        processData.Version,
			CreatedBy:      cmd.SubmittedBy,
			Approver:       *cmd.Approver,
			State:          aopprovaltaskstate.Open,
		}
		createdAt, err := templaterepository.CreateApprovalTask(ctx, tx, data)
		if err != nil {
			return fmt.Errorf("could not create approval task: %w", err)
		}

		createApprovalTaskEvent := templateevents.ContractTemplateCreateApprovalTaskEvent{
			DID:            processData.DID,
			DocumentNumber: processData.DocumentNumber,
			Version:        processData.Version,
			CreatedBy:      cmd.SubmittedBy,
			Approver:       *cmd.Approver,
			OccurredAt:     *createdAt,
		}
		err = event.Create(ctx, tx, createApprovalTaskEvent)
		if err != nil {
			return fmt.Errorf("could not create event: %w", err)
		}

		nextTemplateState = templatestate.Submitted

	} else if processData.State == templatestate.Rejected {

		err := reopenReviewTasks(ctx, tx, cmd.SubmittedBy, processData)
		if err != nil {
			return fmt.Errorf("could not reopen review tasks: %w", err)
		}

		nextTemplateState = templatestate.Submitted

	} else if processData.State == templatestate.Submitted {

		if cmd.ActionFlag != nil {
			if *cmd.ActionFlag == actionflag.Approval {

				err := templaterepository.UpdateReviewTask(ctx, tx, processData.DID, processData.DocumentNumber, processData.Version, cmd.SubmittedBy, reviewtaskstate.Approved)
				if err != nil {
					return fmt.Errorf("could not update approval task: %w", err)
				}

				exist, err := templaterepository.ExistReviewTaskInState(ctx, tx, processData.DID, processData.DocumentNumber, processData.Version, reviewtaskstate.Open)
				if err != nil {
					return fmt.Errorf("could not check if review task exists: %w", err)
				}

				if !exist {
					nextTemplateState = templatestate.Reviewed
				}

			} else if *cmd.ActionFlag == actionflag.Draft {

				err := templaterepository.ReopenReviewTasks(ctx, tx, processData.DID, processData.DocumentNumber, processData.Version)
				if err != nil {
					return fmt.Errorf("could not reopen review tasks: %w", err)
				}

				reopenReviewTaskEvent := templateevents.ContractTemplateReopenReviewTaskEvent{
					DID:            processData.DID,
					DocumentNumber: processData.DocumentNumber,
					Version:        processData.Version,
					CreatedBy:      cmd.SubmittedBy,
					OccurredAt:     time.Now(),
				}
				err = event.Create(ctx, tx, reopenReviewTaskEvent)
				if err != nil {
					return fmt.Errorf("could not create event: %w", err)
				}

				nextTemplateState = templatestate.Rejected
			}
		} else {
			return errors.New("action flags is missing")
		}

	} else if processData.State == templatestate.Reviewed {

		err := reopenReviewTasks(ctx, tx, cmd.SubmittedBy, processData)
		if err != nil {
			return fmt.Errorf("could not reopen review tasks: %w", err)
		}

		nextTemplateState = templatestate.Submitted

	} else {
		return errors.New("current template contract state is invalid")
	}

	if len(nextTemplateState) > 0 && processData.State != nextTemplateState {
		err = templaterepository.UpdateContractTemplateState(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, nextTemplateState)
		if err != nil {
			return fmt.Errorf("could not update contract template state: %w", err)
		}

		evt := templateevents.ContractTemplateSubmittedEvent{
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
