package template_repository

import (
	"context"
	"database/sql"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"errors"
	"fmt"
	"time"
)

type ContractTemplateData struct {
	DID            string                       `db:"did"`
	DocumentNumber int                          `db:"document_number"`
	Version        int                          `db:"version"`
	State          template_state.TemplateState `db:"state"`
	Name           *string                      `db:"name"`
	Description    *string                      `db:"description"`
	CreatedBy      string                       `db:"created_by"`
	CreatedAt      time.Time                    `db:"created_at"`
	UpdatedBy      string                       `db:"updated_by"`
	UpdatedAt      time.Time                    `db:"updated_at"`
	MetaData       *datatype.JSON               `db:"meta_data"`
}

func ReadContractTemplateData(ctx context.Context, tx *sql.Tx, did string) (*ContractTemplateData, error) {
	var ct ContractTemplateData
	err := tx.QueryRowContext(ctx, `
        SELECT did, document_number, version, state, name, description,
               created_by, created_at, updated_by, updated_at, meta_data
        FROM contract_templates WHERE did = $1
    `, did).Scan(
		&ct.DID, &ct.DocumentNumber, &ct.Version, &ct.State, &ct.Name, &ct.Description,
		&ct.CreatedBy, &ct.CreatedAt, &ct.UpdatedBy, &ct.UpdatedAt, &ct.MetaData,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(fmt.Sprintf("contract template with DID %s not found", did))
		}
		return nil, err
	}
	return &ct, nil
}

func ReadContractTemplateState(ctx context.Context, tx *sql.Tx, did string) (*template_state.TemplateState, error) {
	var state template_state.TemplateState
	err := tx.QueryRowContext(ctx, `
        SELECT state
        FROM contract_templates WHERE did = $1
    `, did).Scan(&state)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(fmt.Sprintf("contract template with DID %s not found", did))
		}
		return nil, err
	}
	return &state, nil
}

type ContractTemplateCoreData struct {
	DID            string                       `db:"did"`
	DocumentNumber int                          `db:"document_number"`
	Version        int                          `db:"version"`
	State          template_state.TemplateState `db:"state"`
}

func ReadContractTemplateCoreData(ctx context.Context, tx *sql.Tx, did string) (*ContractTemplateCoreData, error) {
	var coreData ContractTemplateCoreData
	err := tx.QueryRowContext(ctx, `
        SELECT did, document_number, version, state
        FROM contract_templates WHERE did = $1
    `, did).Scan(&coreData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(fmt.Sprintf("contract template with DID %s not found", did))
		}
		return nil, err
	}
	return &coreData, nil
}
