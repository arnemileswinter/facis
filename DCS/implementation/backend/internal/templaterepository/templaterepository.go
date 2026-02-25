package templaterepository

import (
	"context"
	"database/sql"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/approvaltask"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/datatype/templatetype"
	"digital-contracting-service/internal/templaterepository/reviewtask"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type ContractTemplate struct {
	DID            string                      `db:"did"`
	DocumentNumber int                         `db:"document_number"`
	Version        int                         `db:"version"`
	State          templatestate.TemplateState `db:"state"`
	TemplateType   templatetype.TemplateType   `db:"template_type"`
	Name           *string                     `db:"name"`
	Description    *string                     `db:"description"`
	CreatedBy      string                      `db:"created_by"`
	CreatedAt      time.Time                   `db:"created_at"`
	UpdatedAt      time.Time                   `db:"updated_at"`
	TemplateData   *datatype.JSON              `db:"template_data"`
}

func CreateContractTemplate(ctx context.Context, tx *sqlx.Tx, data ContractTemplate) (*time.Time, error) {
	statement := `
    INSERT INTO contract_templates (
        did, created_by, state, name, 
        description, template_data, template_type
    ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING created_at
`

	var createdAt time.Time
	err := tx.GetContext(ctx, &createdAt, statement,
		data.DID,
		data.CreatedBy,
		data.State,
		data.Name,
		data.Description,
		data.TemplateData,
		data.TemplateType,
	)
	if err != nil {
		return nil, err
	}

	return &createdAt, nil
}

////	func GetAmountOfContractTemplatesForID(ctx context.Context, tx *sqlx.Tx, did string, documentNumber *int, version *int) (int, error) {
////		query := `
////	        SELECT COUNT(*)
////	        FROM contract_templates
////	        WHERE did = $1
////		`
////
////		var params []interface{}
////		paramIndex := 1
////
////		if documentNumber != nil {
////			query += `AND document_number = $` + strconv.Itoa(paramIndex) + `,`
////			params = append(params, documentNumber)
////			paramIndex++
////		}
////
////		if version != nil {
////			query += `AND version = $` + strconv.Itoa(paramIndex) + `,`
////			params = append(params, version)
////			paramIndex++
////		}
////
////		return count > 0, nil
////	}

func ReadDataByID(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) (*ContractTemplate, error) {
	query := `
        SELECT did, document_number, version, state, name, description,
               created_by, created_at, updated_at, template_data, template_type
        FROM contract_templates WHERE did = $1 AND document_number = $2 AND version = $3
    `

	var ct ContractTemplate
	err := tx.GetContext(ctx, &ct, query, did, documentNumber, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("contract template with DID %s not found", did)
		}
		return nil, err
	}
	return &ct, nil
}

type MetaData struct {
	DID            string                      `db:"did"`
	DocumentNumber int                         `db:"document_number"`
	Version        int                         `db:"version"`
	State          templatestate.TemplateState `db:"state"`
	TemplateType   templatetype.TemplateType   `db:"template_type"`
	Name           *string                     `db:"name"`
	Description    *string                     `db:"description"`
	CreatedBy      string                      `db:"created_by"`
	CreatedAt      time.Time                   `db:"created_at"`
	UpdatedAt      time.Time                   `db:"updated_at"`
}

func ReadAllMetaData(ctx context.Context, tx *sqlx.Tx) ([]MetaData, error) {
	query := `
        SELECT did, document_number, version, state, template_type, name, description, created_by, created_at, updated_at
        FROM contract_templates
    `

	var cts []MetaData
	err := tx.SelectContext(ctx, &cts, query)
	if err != nil {
		return []MetaData{}, err
	}
	return cts, nil
}

type SearchValues struct {
	DID            *string
	DocumentNumber *int
	Version        *int
	State          *templatestate.TemplateState
	TemplateType   *templatetype.TemplateType
	Name           *string
	Description    *string
	Filter         *string
}

func createSearchConditions(values SearchValues) (*string, []interface{}, error) {

	conditions := ""

	var params []interface{}
	paramIndex := 1

	if values.DID != nil {
		conditions += ` did = $` + strconv.Itoa(paramIndex) + ` AND`
		params = append(params, *values.DID)
		paramIndex++
	}

	if values.DocumentNumber != nil {
		conditions += ` document_number = $` + strconv.Itoa(paramIndex) + ` AND`
		params = append(params, *values.DocumentNumber)
		paramIndex++
	}

	if values.Version != nil {
		conditions += ` version = $` + strconv.Itoa(paramIndex) + ` AND`
		params = append(params, *values.Version)
		paramIndex++
	}

	if values.State != nil {
		conditions += ` state = $` + strconv.Itoa(paramIndex) + ` AND`
		params = append(params, *values.State)
		paramIndex++
	}

	if values.TemplateType != nil {
		conditions += ` template_type = $` + strconv.Itoa(paramIndex) + ` AND`
		params = append(params, "%"+*values.TemplateType+"%")
		paramIndex++
	}

	if values.Name != nil {
		conditions += ` name ILIKE $` + strconv.Itoa(paramIndex) + ` AND`
		params = append(params, "%"+*values.Name+"%")
		paramIndex++
	}

	if values.Description != nil {
		conditions += ` description ILIKE $` + strconv.Itoa(paramIndex) + ` AND`
		params = append(params, "%"+*values.Description+"%")
		paramIndex++
	}

	if values.Filter != nil {
		conditions += ` search_vector @@ plainto_tsquery('english', $` +
			strconv.Itoa(paramIndex) + `) AND`
		params = append(params, *values.Filter)
		paramIndex++
	}

	// Remove last comma " AND"
	l := len(" AND")
	if len(conditions) > l {
		conditions = conditions[:len(conditions)-l]
	}

	return &conditions, params, nil
}

func ReadAllMetaDataByFilter(ctx context.Context, tx *sqlx.Tx, values SearchValues) ([]MetaData, error) {
	query := `
SELECT did, document_number, version, state, name, template_type, description, created_by, created_at, updated_at
FROM contract_templates
`

	conditions, params, err := createSearchConditions(values)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		query += " WHERE " + *conditions
	}

	var cts []MetaData
	_ = cts
	err = tx.SelectContext(ctx, &cts, query, params...)
	if err != nil {
		return []MetaData{}, err
	}
	return cts, nil
}

type ProcessData struct {
	DID            string                      `db:"did"`
	DocumentNumber int                         `db:"document_number"`
	Version        int                         `db:"version"`
	State          templatestate.TemplateState `db:"state"`
	CreatedBy      string                      `db:"created_by"`
	UpdatedAt      time.Time                   `db:"updated_at"`
}

func ReadProcessData(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) (*ProcessData, error) {
	query := `
        SELECT did, document_number, version, state, updated_at, created_by
        FROM contract_templates WHERE did = $1 AND document_number = $2 AND version = $3
    `

	var processData ProcessData
	err := tx.GetContext(ctx, &processData, query, did, documentNumber, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("contract template with DID %s, DocumentNumber %d and Version %d not found", did, documentNumber, version)
		}
		return nil, err
	}
	return &processData, nil
}

func UpdateState(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int, state templatestate.TemplateState) error {
	statement := `UPDATE contract_templates SET
        	state = $4
    	WHERE did = $1 AND document_number = $2 AND version = $3	
`
	_, err := tx.ExecContext(ctx, statement, did, documentNumber, version, state)
	if err != nil {
		return err
	}

	return err
}

func createQuery(data ContractTemplateUpdateData) (*string, []interface{}, error) {

	queryBase := `UPDATE contract_templates SET `

	var columns []string
	var params []interface{}

	addParam := func(columnName string, value interface{}) {
		columns = append(columns, fmt.Sprintf("%s = $%d", columnName, len(params)+1))
		params = append(params, value)
	}

	if data.Name != nil {
		addParam("name", data.Name)
	}

	if data.Description != nil {
		addParam("description", data.Description)
	}

	if data.TemplateData != nil && data.TemplateData.IsNotNullValue() {
		addParam("template_data", data.TemplateData)
	}

	if data.TemplateType != nil && data.TemplateData.IsNotNullValue() {
		addParam("template_type", data.TemplateType)
	}

	if len(columns) == 0 {
		return nil, nil, errors.New("no fields to update")
	}

	fullQuery := queryBase + strings.Join(columns, ", ")

	nextIdx := len(params) + 1
	fullQuery += fmt.Sprintf(" WHERE did = $%d AND document_number = $%d AND version = $%d;",
		nextIdx, nextIdx+1, nextIdx+2)

	params = append(params, data.DID, data.DocumentNumber, data.Version)

	return &fullQuery, params, nil
}

type ContractTemplateUpdateData struct {
	DID            string `db:"did"`
	DocumentNumber int    `db:"document_number"`
	Version        int    `db:"version"`

	TemplateType *templatetype.TemplateType `db:"template_type"`
	Name         *string                    `db:"name"`
	Description  *string                    `db:"description"`
	TemplateData *datatype.JSON             `db:"template_data"`
}

func UpdateData(ctx context.Context, tx *sqlx.Tx, data ContractTemplateUpdateData) error {
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

func ReopenTasks(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) error {
	err := reviewtask.ReopenTasks(ctx, tx, did, documentNumber, version)
	if err != nil {
		return fmt.Errorf("could not reopen review tasks: %w", err)
	}

	err = approvaltask.ReopenTasks(ctx, tx, did, documentNumber, version)
	if err != nil {
		return fmt.Errorf("could not reopen approval tasks: %w", err)
	}

	return nil
}

func CleanupTasks(ctx context.Context, tx *sqlx.Tx, did string, documentNumber int, version int) error {
	err := reviewtask.DeleteTask(ctx, tx, did, documentNumber, version)
	if err != nil {
		return fmt.Errorf("could not delete review task: %w", err)
	}

	err = approvaltask.DeleteTask(ctx, tx, did, documentNumber, version)
	if err != nil {
		return fmt.Errorf("could not delete approval task: %w", err)
	}

	return nil
}
