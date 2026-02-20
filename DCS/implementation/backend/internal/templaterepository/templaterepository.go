package templaterepository

import (
	"context"
	"database/sql"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type ContractTemplateData struct {
	DID            string                      `db:"did"`
	DocumentNumber int                         `db:"document_number"`
	Version        int                         `db:"version"`
	State          templatestate.TemplateState `db:"state"`
	Name           *string                     `db:"name"`
	Description    *string                     `db:"description"`
	CreatedBy      string                      `db:"created_by"`
	CreatedAt      time.Time                   `db:"created_at"`
	UpdatedAt      time.Time                   `db:"updated_at"`
	TemplateData   *datatype.JSON              `db:"template_data"`
}

func CreateContractTemplate(ctx context.Context, tx *sqlx.Tx, data ContractTemplateData) (*time.Time, error) {
	query := `
    INSERT INTO contract_templates (
        did, created_by, state, name, 
        description, template_data
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
		data.TemplateData,
	)
	if err != nil {
		return nil, err
	}

	return &createdAt, nil
}

func ReadContractTemplateDataById(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) (*ContractTemplateData, error) {
	query := `
        SELECT did, document_number, version, state, name, description,
               created_by, created_at, updated_at, template_data
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

type ContractTemplateMetaData struct {
	DID            string                      `db:"did"`
	DocumentNumber int                         `db:"document_number"`
	Version        int                         `db:"version"`
	State          templatestate.TemplateState `db:"state"`
	Name           *string                     `db:"name"`
	Description    *string                     `db:"description"`
	CreatedBy      string                      `db:"created_by"`
	CreatedAt      time.Time                   `db:"created_at"`
	UpdatedAt      time.Time                   `db:"updated_at"`
}

func ReadAllContractTemplateMetaData(ctx context.Context, tx *sqlx.Tx) ([]ContractTemplateMetaData, error) {
	query := `
        SELECT did, document_number, version, state, name, description, created_by, created_at, updated_at
        FROM contract_templates
    `

	var cts []ContractTemplateMetaData
	err := tx.SelectContext(ctx, &cts, query)
	if err != nil {
		return []ContractTemplateMetaData{}, err
	}
	return cts, nil
}

/*
	func createSearchQuery(sortedKey []string, filter map[string]interface{}) (*string, []interface{}, error) {
		query := `

SELECT did, document_number, version, state, name, description, created_by, created_at, updated_at
FROM contract_templates
WHERE
`

		var params []interface{}
		paramIndex := 1

		for i, k := range sortedKey {

		}

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

		if data.TemplateData != nil && data.TemplateData.IsNotNullValue() {
			query += ` template_data = $` + strconv.Itoa(paramIndex) + `,`
			params = append(params, data.TemplateData)
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

	func ReadAllContractTemplateMetaDataByFilter(ctx context.Context, tx *sqlx.Tx, filter map[string]interface{}) ([]ContractTemplateMetaData, error) {
		query := `
		    SELECT did, document_number, version, state, name, description, created_at, updated_at, meta_data
		    FROM contract_templates
		    WHERE state = $1
		 `

		keys := make([]string, 0, len(filter))
		for k := range filter {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		var cts []ContractTemplateMetaData
		_ = cts
		err := tx.SelectContext(ctx, &cts, query, state)
		if err != nil {
			return []ContractTemplateMetaData{}, err
		}
		return cts, nil
	}
*/
type ContractTemplateProcessData struct {
	DID            string                      `db:"did"`
	DocumentNumber int                         `db:"document_number"`
	Version        int                         `db:"version"`
	State          templatestate.TemplateState `db:"state"`
	CreatedBy      string                      `db:"created_by"`
	UpdatedAt      time.Time                   `db:"updated_at"`
}

func ReadContractTemplateProcessData(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) (*ContractTemplateProcessData, error) {
	query := `
        SELECT did, document_number, version, state, updated_at, created_by
        FROM contract_templates WHERE did = $1 AND document_number = $2 AND version = $3
    `

	var processData ContractTemplateProcessData
	err := tx.GetContext(ctx, &processData, query, did, documentNumber, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(fmt.Sprintf("contract template with DID %s not found", did))
		}
		return nil, err
	}
	return &processData, nil
}

func UpdateContractTemplateState(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int, state templatestate.TemplateState) error {
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

	if data.TemplateData != nil && data.TemplateData.IsNotNullValue() {
		query += ` template_data = $` + strconv.Itoa(paramIndex) + `,`
		params = append(params, data.TemplateData)
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
