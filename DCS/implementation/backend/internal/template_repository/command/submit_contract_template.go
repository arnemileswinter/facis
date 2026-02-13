package command

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/template_repository"
	"digital-contracting-service/internal/template_repository/datatype/action_flag"
	aopprovaltaskstate "digital-contracting-service/internal/template_repository/datatype/approval_task_state"
	"digital-contracting-service/internal/template_repository/datatype/review_task_state"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	templateevents "digital-contracting-service/internal/template_repository/event"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type SubmitContractTemplateCommand struct {
	DID            string
	DocumentNumber int
	Version        int
	SubmittedBy    string
	ActionFlag     *action_flag.ActionFlag
	Comments       []string
	Reviewer       []string
	Approver       *string
}

type SubmitContractTemplateHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *SubmitContractTemplateHandler) Handle(cmd SubmitContractTemplateCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.GetTransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	coreData, err := template_repository.ReadContractTemplateCoreData(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return err
	}

	var nextTemplateState template_state.TemplateState
	if coreData.State == template_state.Draft {

		if cmd.Reviewer == nil || len(cmd.Reviewer) == 0 {
			return errors.New("no reviewer provided")
		}

		if cmd.Approver == nil || len(*cmd.Approver) == 0 {
			return errors.New("no approver provided")
		}

		for _, reviewer := range cmd.Reviewer {
			reviewTask := template_repository.ReviewTaskData{
				DID:            cmd.DID,
				DocumentNumber: coreData.DocumentNumber,
				Version:        coreData.Version,
				Reviewer:       reviewer,
				State:          review_task_state.Open,
				CreatedBy:      cmd.SubmittedBy,
			}
			createdAt, err := template_repository.CreateReviewTasks(ctx, tx, reviewTask)
			if err != nil {
				return err
			}

			createReviewTaskEvent := templateevents.ContractTemplateCreateReviewTaskEvent{
				DID:            coreData.DID,
				DocumentNumber: coreData.DocumentNumber,
				Version:        coreData.Version,
				CreatedBy:      cmd.SubmittedBy,
				Reviewer:       reviewer,
				OccurredAt:     *createdAt,
			}
			err = event.CreateNewEvent(ctx, tx, createReviewTaskEvent)
			if err != nil {
				return err
			}
		}

		data := template_repository.ApprovalTaskData{
			DID:            cmd.DID,
			DocumentNumber: coreData.DocumentNumber,
			Version:        coreData.Version,
			CreatedBy:      coreData.CreatedBy,
			Approver:       *cmd.Approver,
			State:          aopprovaltaskstate.Open,
		}
		createdAt, err := template_repository.CreateApprovalTask(ctx, tx, data)
		if err != nil {
			return err
		}

		createApprovalTaskEvent := templateevents.ContractTemplateCreateApprovalTaskEvent{
			DID:            coreData.DID,
			DocumentNumber: coreData.DocumentNumber,
			Version:        coreData.Version,
			CreatedBy:      cmd.SubmittedBy,
			Approver:       *cmd.Approver,
			OccurredAt:     *createdAt,
		}
		err = event.CreateNewEvent(ctx, tx, createApprovalTaskEvent)
		if err != nil {
			return err
		}

		nextTemplateState = template_state.Submitted

	} else if coreData.State == template_state.Rejected {

		err := template_repository.ReopenReviewTasks(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version)
		if err != nil {
			return err
		}

		reopenReviewTaskEvent := templateevents.ContractTemplateReopenReviewTaskEvent{
			DID:            coreData.DID,
			DocumentNumber: coreData.DocumentNumber,
			Version:        coreData.Version,
			CreatedBy:      cmd.SubmittedBy,
			OccurredAt:     time.Now(),
		}
		err = event.CreateNewEvent(ctx, tx, reopenReviewTaskEvent)
		if err != nil {
			return err
		}

		err = template_repository.ReopenApprovalTask(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version)
		if err != nil {
			return err
		}

		reopenApprovalTaskEvent := templateevents.ContractTemplateReopenApprovalTaskEvent{
			DID:            coreData.DID,
			DocumentNumber: coreData.DocumentNumber,
			Version:        coreData.Version,
			CreatedBy:      cmd.SubmittedBy,
			OccurredAt:     time.Now(),
		}
		err = event.CreateNewEvent(ctx, tx, reopenApprovalTaskEvent)
		if err != nil {
			return err
		}

		nextTemplateState = template_state.Submitted

	} else if coreData.State == template_state.Submitted {

		if cmd.ActionFlag != nil {
			if *cmd.ActionFlag == action_flag.Approval {

				err := template_repository.UpdateReviewTask(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version, cmd.SubmittedBy, review_task_state.Approved)
				if err != nil {
					return err
				}

				exist, err := template_repository.ExistReviewTaskInState(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version, review_task_state.Open)
				if err != nil {
					return err
				}

				if !exist {
					nextTemplateState = template_state.Reviewed
				}

			} else if *cmd.ActionFlag == action_flag.Draft {

				err := template_repository.ReopenReviewTasks(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version)
				if err != nil {
					return err
				}

				reopenReviewTaskEvent := templateevents.ContractTemplateReopenReviewTaskEvent{
					DID:            coreData.DID,
					DocumentNumber: coreData.DocumentNumber,
					Version:        coreData.Version,
					CreatedBy:      cmd.SubmittedBy,
					OccurredAt:     time.Now(),
				}
				err = event.CreateNewEvent(ctx, tx, reopenReviewTaskEvent)
				if err != nil {
					return err
				}

				nextTemplateState = template_state.Rejected
			}
		} else {
			return errors.New("action flags is missing")
		}

	} else if coreData.State == template_state.Reviewed {

		err := template_repository.ReopenReviewTasks(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version)
		if err != nil {
			return err
		}

		reopenReviewTaskEvent := templateevents.ContractTemplateReopenReviewTaskEvent{
			DID:            coreData.DID,
			DocumentNumber: coreData.DocumentNumber,
			Version:        coreData.Version,
			CreatedBy:      cmd.SubmittedBy,
			OccurredAt:     time.Now(),
		}
		err = event.CreateNewEvent(ctx, tx, reopenReviewTaskEvent)
		if err != nil {
			return err
		}

		err = template_repository.ReopenApprovalTask(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version)
		if err != nil {
			return err
		}

		reopenApprovalTaskEvent := templateevents.ContractTemplateReopenApprovalTaskEvent{
			DID:            coreData.DID,
			DocumentNumber: coreData.DocumentNumber,
			Version:        coreData.Version,
			CreatedBy:      cmd.SubmittedBy,
			OccurredAt:     time.Now(),
		}
		err = event.CreateNewEvent(ctx, tx, reopenApprovalTaskEvent)
		if err != nil {
			return err
		}

		nextTemplateState = template_state.Submitted

	} else {
		return errors.New("current template contract state is invalid")
	}

	if len(nextTemplateState) > 0 && coreData.State != nextTemplateState {
		err = template_repository.UpdateContractTemplateState(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, nextTemplateState)
		if err != nil {
			return err
		}

		evt := templateevents.ContractTemplateSubmittedEvent{
			DID:            cmd.DID,
			DocumentNumber: cmd.DocumentNumber,
			Version:        cmd.DocumentNumber,
			SubmittedBy:    cmd.SubmittedBy,
			PreviousState:  coreData.State,
			NewState:       nextTemplateState,
			ActionFlag:     cmd.ActionFlag,
			Comments:       cmd.Comments,
			OccurredAt:     time.Now(),
		}
		err = event.CreateNewEvent(ctx, tx, evt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
