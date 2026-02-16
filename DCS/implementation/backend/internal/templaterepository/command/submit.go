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
	"time"

	"github.com/jmoiron/sqlx"
)

type SubmitContractTemplateCommand struct {
	DID            string
	DocumentNumber int
	Version        int
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

func reopenReviewTasks(ctx context.Context, tx *sqlx.Tx, submittedBy string, data *templaterepository.ContractTemplateCoreData) error {
	err := templaterepository.ReopenReviewTasks(ctx, tx, data.DID, data.DocumentNumber, data.Version)
	if err != nil {
		return err
	}

	reopenReviewTaskEvent := templateevents.ContractTemplateReopenReviewTaskEvent{
		DID:            data.DID,
		DocumentNumber: data.DocumentNumber,
		Version:        data.Version,
		CreatedBy:      submittedBy,
		OccurredAt:     time.Now(),
	}
	err = event.CreateNewEvent(ctx, tx, reopenReviewTaskEvent)
	if err != nil {
		return err
	}

	err = templaterepository.ReopenApprovalTask(ctx, tx, data.DID, data.DocumentNumber, data.Version)
	if err != nil {
		return err
	}

	reopenApprovalTaskEvent := templateevents.ContractTemplateReopenApprovalTaskEvent{
		DID:            data.DID,
		DocumentNumber: data.DocumentNumber,
		Version:        data.Version,
		CreatedBy:      submittedBy,
		OccurredAt:     time.Now(),
	}
	err = event.CreateNewEvent(ctx, tx, reopenApprovalTaskEvent)
	if err != nil {
		return err
	}

	return nil
}

func (h *SubmitContractTemplateHandler) Handle(cmd SubmitContractTemplateCommand) error {

	ctx, cancel := context.WithTimeout(h.Ctx, base.GetTransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	coreData, err := templaterepository.ReadContractTemplateCoreData(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version)
	if err != nil {
		return err
	}

	var nextTemplateState templatestate.TemplateState
	if coreData.State == templatestate.Draft {

		if cmd.Reviewer == nil || len(cmd.Reviewer) == 0 {
			return errors.New("no reviewer provided")
		}

		if cmd.Approver == nil || len(*cmd.Approver) == 0 {
			return errors.New("no approver provided")
		}

		for _, reviewer := range cmd.Reviewer {
			reviewTask := templaterepository.ReviewTaskData{
				DID:            cmd.DID,
				DocumentNumber: coreData.DocumentNumber,
				Version:        coreData.Version,
				Reviewer:       reviewer,
				State:          reviewtaskstate.Open,
				CreatedBy:      cmd.SubmittedBy,
			}
			createdAt, err := templaterepository.CreateReviewTasks(ctx, tx, reviewTask)
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

		data := templaterepository.ApprovalTaskData{
			DID:            cmd.DID,
			DocumentNumber: coreData.DocumentNumber,
			Version:        coreData.Version,
			CreatedBy:      coreData.CreatedBy,
			Approver:       *cmd.Approver,
			State:          aopprovaltaskstate.Open,
		}
		createdAt, err := templaterepository.CreateApprovalTask(ctx, tx, data)
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

		nextTemplateState = templatestate.Submitted

	} else if coreData.State == templatestate.Rejected {

		err := reopenReviewTasks(ctx, tx, cmd.SubmittedBy, coreData)
		if err != nil {
			return err
		}

		nextTemplateState = templatestate.Submitted

	} else if coreData.State == templatestate.Submitted {

		if cmd.ActionFlag != nil {
			if *cmd.ActionFlag == actionflag.Approval {

				err := templaterepository.UpdateReviewTask(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version, cmd.SubmittedBy, reviewtaskstate.Approved)
				if err != nil {
					return err
				}

				exist, err := templaterepository.ExistReviewTaskInState(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version, reviewtaskstate.Open)
				if err != nil {
					return err
				}

				if !exist {
					nextTemplateState = templatestate.Reviewed
				}

			} else if *cmd.ActionFlag == actionflag.Draft {

				err := templaterepository.ReopenReviewTasks(ctx, tx, coreData.DID, coreData.DocumentNumber, coreData.Version)
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

				nextTemplateState = templatestate.Rejected
			}
		} else {
			return errors.New("action flags is missing")
		}

	} else if coreData.State == templatestate.Reviewed {

		err := reopenReviewTasks(ctx, tx, cmd.SubmittedBy, coreData)
		if err != nil {
			return err
		}

		nextTemplateState = templatestate.Submitted

	} else {
		return errors.New("current template contract state is invalid")
	}

	if len(nextTemplateState) > 0 && coreData.State != nextTemplateState {
		err = templaterepository.UpdateContractTemplateState(ctx, tx, cmd.DID, cmd.DocumentNumber, cmd.Version, nextTemplateState)
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
