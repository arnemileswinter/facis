package command

import (
	"context"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/template_repository"
	"digital-contracting-service/internal/template_repository/datatype/action_flag"
	"digital-contracting-service/internal/template_repository/datatype/review_task_state"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	templateevents "digital-contracting-service/internal/template_repository/event"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type SubmitTemplateContractCommand struct {
	DID            string
	SubmittedBy    string
	ActionFlag     *action_flag.ActionFlag
	ReviewComments []string
	Assignees      []string
}

type SubmitTemplateContractHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func (h *SubmitTemplateContractHandler) Handle(cmd SubmitTemplateContractCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	coreData, err := template_repository.ReadContractTemplateCoreData(ctx, tx, cmd.DID)
	if err != nil {
		return err
	}

	var nextTemplateState template_state.TemplateState
	if coreData.State == template_state.Draft || coreData.State == template_state.Rejected {

		for _, assignee := range cmd.Assignees {
			reviewTask := template_repository.ReviewTaskData{
				DID:            cmd.DID,
				DocumentNumber: coreData.DocumentNumber,
				Version:        coreData.Version,
				Assignee:       assignee,
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
				Assignee:       assignee,
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

				err := template_repository.UpdateReviewTask(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version, cmd.SubmittedBy, review_task_state.Approved, cmd.ReviewComments)
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

				err := template_repository.UpdateReviewTask(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version, cmd.SubmittedBy, review_task_state.Rejected, cmd.ReviewComments)
				if err != nil {
					return err
				}

				err = template_repository.CancelReviewTasks(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version)
				if err != nil {
					return err
				}

				nextTemplateState = template_state.Rejected
			}
		} else {
			return errors.New("action flags is missing")
		}

	} else if coreData.State == template_state.Reviewed {

		nextTemplateState = template_state.Approved

	} else {
		return errors.New("current template contract state is invalid")
	}

	if len(nextTemplateState) > 0 && coreData.State != nextTemplateState {
		err = template_repository.UpdateContractTemplateState(ctx, tx, cmd.DID, nextTemplateState)
		if err != nil {
			return err
		}

		evt := templateevents.ContractTemplateSubmittedEvent{
			DID:            cmd.DID,
			SubmittedBy:    cmd.SubmittedBy,
			PreviousState:  coreData.State,
			NewState:       nextTemplateState,
			ActionFlag:     cmd.ActionFlag,
			ReviewComments: cmd.ReviewComments,
			OccurredAt:     time.Now(),
		}
		err = event.CreateNewEvent(ctx, tx, evt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
