package command

import (
	"context"
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

	ctx, cancel := context.WithTimeout(h.Ctx, 5*time.Second)
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

	var approver string
	if cmd.Approver == nil || len(*cmd.Approver) == 0 {
		if coreData.Approver == nil || len(*coreData.Approver) == 0 {
			return errors.New("no approver provided")
		}

		approver = *coreData.Approver
	} else {
		approver = *cmd.Approver
	}

	var nextTemplateState template_state.TemplateState
	if coreData.State == template_state.Draft || coreData.State == template_state.Rejected {

		err := template_repository.UpdateContractTemplateApprover(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, approver)
		if err != nil {
			return err
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
			createdAt, err := template_repository.CreateReviewTask(ctx, tx, reviewTask)
			if err != nil {
				return err
			}

			evt := templateevents.ContractTemplateCreateReviewTaskEvent{
				DID:            coreData.DID,
				DocumentNumber: coreData.DocumentNumber,
				Version:        coreData.Version,
				CreatedBy:      cmd.SubmittedBy,
				Reviewer:       reviewer,
				OccurredAt:     *createdAt,
			}
			err = event.CreateNewEvent(ctx, tx, evt)
			if err != nil {
				return err
			}
		}

		nextTemplateState = template_state.Submitted

	} else if coreData.State == template_state.Submitted {

		if cmd.ActionFlag != nil {
			if *cmd.ActionFlag == action_flag.Approval {

				err := template_repository.UpdateReviewTask(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version, cmd.SubmittedBy, review_task_state.Approved, cmd.Comments)
				if err != nil {
					return err
				}

				exist, err := template_repository.ExistReviewTaskInState(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version, review_task_state.Open)
				if err != nil {
					return err
				}

				if !exist {
					data := template_repository.ApprovalTaskData{
						DID:            cmd.DID,
						DocumentNumber: coreData.DocumentNumber,
						Version:        coreData.Version,
						CreatedBy:      coreData.CreatedBy,
						Approver:       *coreData.Approver,
						State:          aopprovaltaskstate.Open,
					}
					_, err = template_repository.CreateApprovalTask(ctx, tx, data)
					if err != nil {
						return err
					}

					nextTemplateState = template_state.Reviewed
				}

			} else if *cmd.ActionFlag == action_flag.Draft {

				err := template_repository.UpdateReviewTask(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version, cmd.SubmittedBy, review_task_state.Rejected, cmd.Comments)
				if err != nil {
					return err
				}

				err = template_repository.CancelOldReviewTasks(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version)
				if err != nil {
					return err
				}

				nextTemplateState = template_state.Rejected
			}
		} else {
			return errors.New("action flags is missing")
		}

	} else if coreData.State == template_state.Reviewed {

		err := template_repository.UpdateApprovalTask(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version, cmd.SubmittedBy, aopprovaltaskstate.Resubmitted, cmd.Comments)
		if err != nil {
			return err
		}

		err = template_repository.CreateResubmissionTasks(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version, cmd.SubmittedBy)
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
