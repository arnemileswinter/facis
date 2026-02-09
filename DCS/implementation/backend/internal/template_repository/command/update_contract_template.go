package command

import (
	"context"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	"digital-contracting-service/internal/template_repository"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	templateevents "digital-contracting-service/internal/template_repository/event"
	"errors"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type UpdateTemplateContractCommand struct {
	DID         string
	UpdatedBy   string
	Name        *string
	Description *string
	MetaData    *datatype.JSON
}

type UpdateTemplateContractHandler struct {
	Ctx context.Context
	DB  *sqlx.DB
}

func createQuery(cmd UpdateTemplateContractCommand) (*string, []interface{}, error) {
	query := `UPDATE contract_templates SET`

	var params []interface{}
	paramIndex := 1

	if cmd.Name != nil {
		query += ` name = $` + strconv.Itoa(paramIndex) + `,`
		params = append(params, cmd.Name)
		paramIndex++
	}

	if cmd.Description != nil {
		query += ` description = $` + strconv.Itoa(paramIndex) + `,`
		params = append(params, cmd.Description)
		paramIndex++
	}

	if cmd.MetaData != nil && cmd.MetaData.IsNotNullValue() {
		query += ` meta_data = $` + strconv.Itoa(paramIndex) + `,`
		params = append(params, cmd.MetaData)
		paramIndex++
	}

	if len(params) == 0 {
		return nil, nil, errors.New("no parameters found")
	}

	// Remove last comma
	query = query[:len(query)-1]

	query += ` WHERE did = $` + strconv.Itoa(paramIndex) + `;`
	params = append(params, cmd.DID)

	return &query, params, nil
}

func (h *UpdateTemplateContractHandler) Handle(cmd UpdateTemplateContractCommand) error {

	query, params, err := createQuery(cmd)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(h.Ctx, 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}

	oldData, err := template_repository.ReadContractTemplateData(ctx, tx, cmd.DID)
	if err != nil {
		return err
	}

	if oldData.State != template_state.Draft {
		return errors.New("invalid contract template state")
	}

	_, err = tx.ExecContext(ctx, *query, params...)
	if err != nil {
		return err
	}

	evt := templateevents.ContractTemplateUpdatedEvent{
		DID:            cmd.DID,
		UpdatedBy:      cmd.UpdatedBy,
		OldName:        oldData.Name,
		NewName:        cmd.Name,
		OldDescription: oldData.Description,
		NewDescription: cmd.Description,
		OldMetaData:    oldData.MetaData,
		NewMetaData:    cmd.MetaData,
		OccurredAt:     time.Now(),
	}
	err = event.CreateNewEvent(h.Ctx, tx, evt)
	if err != nil {
		return err
	}

	return tx.Commit()
}
