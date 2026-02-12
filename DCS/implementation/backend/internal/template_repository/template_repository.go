package template_repository

import (
	"context"
	"database/sql"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type ContractTemplateData struct {
	DID            string                       `db:"did"`
	DocumentNumber int                          `db:"document_number"`
	Version        int                          `db:"version"`
	State          template_state.TemplateState `db:"state"`
	Name           *string                      `db:"name"`
	Description    *string                      `db:"description"`
	Approver       *string                      `db:"approver"`
	CreatedBy      string                       `db:"created_by"`
	CreatedAt      time.Time                    `db:"created_at"`
	MetaData       *datatype.JSON               `db:"meta_data"`
}

func CreateContractTemplate(ctx context.Context, tx *sqlx.Tx, data ContractTemplateData) (*time.Time, error) {
	query := `
    INSERT INTO contract_templates (
        did, created_by, state, name, 
        description, meta_data
    ) VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING created_at
`

	var createdAt time.Time
	err := tx.GetContext(ctx, &createdAt, query,
		data.DID,
		data.CreatedBy,
		data.State,
		data.Name,
		data.Description,
		data.MetaData,
	)
	if err != nil {
		return nil, err
	}

	return &createdAt, nil
}

func ReadContractTemplateData(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) (*ContractTemplateData, error) {
	query := `
        SELECT did, document_number, version, state, name, description,
               created_by, created_at, meta_data
        FROM contract_templates WHERE did = $1 AND document_number = $2 AND version = $3
    `

	var ct ContractTemplateData
	err := tx.GetContext(ctx, &ct, query, did, documentNumber, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(fmt.Sprintf("contract template with DID %s not found", did))
		}
		return nil, err
	}
	return &ct, nil
}

func ReadAllContractTemplateData(ctx context.Context, tx *sqlx.Tx) ([]ContractTemplateData, error) {
	query := `
        SELECT did, document_number, version, state, name, description,
               created_by, created_at, meta_data
        FROM contract_templates
    `

	var cts []ContractTemplateData
	err := tx.SelectContext(ctx, &cts, query)
	if err != nil {
		return []ContractTemplateData{}, err
	}
	return cts, nil
}

type ContractTemplateCoreData struct {
	DID            string                       `db:"did"`
	DocumentNumber int                          `db:"document_number"`
	Version        int                          `db:"version"`
	State          template_state.TemplateState `db:"state"`
	CreatedBy      string                       `db:"created_by"`
}

func ReadContractTemplateCoreData(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) (*ContractTemplateCoreData, error) {
	query := `
        SELECT did, document_number, version, state, created_by
        FROM contract_templates WHERE did = $1 AND document_number = $2 AND version = $3
    `

	var coreData ContractTemplateCoreData
	err := tx.GetContext(ctx, &coreData, query, did, documentNumber, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(fmt.Sprintf("contract template with DID %s not found", did))
		}
		return nil, err
	}
	return &coreData, nil
}

func ReadContractTemplateState(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) (*template_state.TemplateState, error) {
	query := `
        SELECT state
        FROM contract_templates WHERE did = $1 AND document_number = $2 AND version = $3
    `

	var state template_state.TemplateState
	err := tx.GetContext(ctx, &state, query, did, documentNumber, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(fmt.Sprintf("contract template with DID %s not found", did))
		}
		return nil, err
	}
	return &state, nil
}

func UpdateContractTemplateState(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int, state template_state.TemplateState) error {
	query := `UPDATE contract_templates SET
        	state = $4
    	WHERE did = $1 AND document_number = $2 AND version = $3	
`
	_, err := tx.ExecContext(ctx, query, did, documentNumber, version, state)
	if err != nil {
		return err
	}

	return err
}

func createQuery(data ContractTemplateData) (*string, []interface{}, error) {
	query := `UPDATE contract_templates SET`

	var params []interface{}
	paramIndex := 1

	if data.Name != nil {
		query += ` name = $` + strconv.Itoa(paramIndex) + `,`
		params = append(params, data.Name)
		paramIndex++
	}

	if data.Description != nil {
		query += ` description = $` + strconv.Itoa(paramIndex) + `,`
		params = append(params, data.Description)
		paramIndex++
	}

	if data.MetaData != nil && data.MetaData.IsNotNullValue() {
		query += ` meta_data = $` + strconv.Itoa(paramIndex) + `,`
		params = append(params, data.MetaData)
		paramIndex++
	}

	if len(params) == 0 {
		return nil, nil, errors.New("no parameters found")
	}

	// Remove last comma
	query = query[:len(query)-1]

	query += ` WHERE did = $` + strconv.Itoa(paramIndex) + ` AND document_number = $` + strconv.Itoa(paramIndex+1) + ` AND version = $` + strconv.Itoa(paramIndex+2) + `;`
	params = append(params, data.DID, data.DocumentNumber, data.Version)

	return &query, params, nil
}

func UpdateTemplateContractData(ctx context.Context, tx *sqlx.Tx, data ContractTemplateData) error {
	query, params, err := createQuery(data)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, *query, params...)
	if err != nil {
		return err
	}

	return err
}
