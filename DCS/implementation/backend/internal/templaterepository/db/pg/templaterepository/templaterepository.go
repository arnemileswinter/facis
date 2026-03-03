package templaterepository

import (
	"context"
	"database/sql"
	"digital-contracting-service/internal/templaterepository/datatype/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type PostgresContractTemplateRepo struct {
	Ctx context.Context
}

func (r *PostgresContractTemplateRepo) Create(tx *sqlx.Tx, data templaterepository.ContractTemplate) (*time.Time, error) {
	statement := `
        INSERT INTO contract_templates (
            did, created_by, state, name,
            description, template_data, template_type
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING created_at
    `
	var createdAt time.Time
	err := tx.GetContext(r.Ctx, &createdAt, statement,
		data.DID, data.CreatedBy, data.State, data.Name,
		data.Description, data.TemplateData, data.TemplateType,
	)
	if err != nil {
		return nil, err
	}
	return &createdAt, nil
}

func (r *PostgresContractTemplateRepo) ReadDataByID(tx *sqlx.Tx, did string, documentNumber int, version int) (*templaterepository.ContractTemplate, error) {
	query := `
        SELECT did, document_number, version, state, name, description,
               created_by, created_at, updated_at, template_data, template_type
        FROM contract_templates WHERE did = $1 AND document_number = $2 AND version = $3
    `
	var ct templaterepository.ContractTemplate
	err := tx.GetContext(r.Ctx, &ct, query, did, documentNumber, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("contract template with DID %s not found", did)
		}
		return nil, err
	}
	return &ct, nil
}

func (r *PostgresContractTemplateRepo) ReadAllMetaData(tx *sqlx.Tx) ([]templaterepository.MetaData, error) {
	query := `
        SELECT did, document_number, version, state, template_type, name, description, created_by, created_at, updated_at
        FROM contract_templates
    `
	var cts []templaterepository.MetaData
	err := tx.SelectContext(r.Ctx, &cts, query)
	if err != nil {
		return []templaterepository.MetaData{}, err
	}
	return cts, nil
}

func (r *PostgresContractTemplateRepo) ReadAllMetaDataByFilter(tx *sqlx.Tx, values templaterepository.SearchValues) ([]templaterepository.MetaData, error) {
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

	var cts []templaterepository.MetaData
	err = tx.SelectContext(r.Ctx, &cts, query, params...)
	if err != nil {
		return []templaterepository.MetaData{}, err
	}
	return cts, nil
}

func (r *PostgresContractTemplateRepo) ReadProcessData(tx *sqlx.Tx, did string, documentNumber int, version int) (*templaterepository.ProcessData, error) {
	query := `
        SELECT did, document_number, version, state, updated_at, created_by
        FROM contract_templates WHERE did = $1 AND document_number = $2 AND version = $3
    `
	var processData templaterepository.ProcessData
	err := tx.GetContext(r.Ctx, &processData, query, did, documentNumber, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("contract template with DID %s, DocumentNumber %d and Version %d not found", did, documentNumber, version)
		}
		return nil, err
	}
	return &processData, nil
}

func (r *PostgresContractTemplateRepo) UpdateState(tx *sqlx.Tx, did string, documentNumber int, version int, state templatestate.TemplateState) error {
	statement := `
        UPDATE contract_templates SET state = $4
        WHERE did = $1 AND document_number = $2 AND version = $3
    `
	_, err := tx.ExecContext(r.Ctx, statement, did, documentNumber, version, state)
	return err
}

func (r *PostgresContractTemplateRepo) Update(tx *sqlx.Tx, data templaterepository.UpdateData) error {
	query, params, err := createQuery(data)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(r.Ctx, *query, params...)
	return err
}

// --- Hilfsfunktionen (package-intern, keine Methoden) ---

func createSearchConditions(values templaterepository.SearchValues) (*string, []interface{}, error) {
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
		conditions += ` search_vector @@ plainto_tsquery('english', $` + strconv.Itoa(paramIndex) + `) AND`
		params = append(params, *values.Filter)
		paramIndex++
	}

	l := len(" AND")
	if len(conditions) > l {
		conditions = conditions[:len(conditions)-l]
	}

	return &conditions, params, nil
}

func createQuery(data templaterepository.UpdateData) (*string, []interface{}, error) {
	queryBase := `UPDATE contract_templates SET `
	var columns []string
	var params []interface{}

	addParam := func(columnName string, value interface{}) {
		columns = append(columns, fmt.Sprintf("%s = $%d", columnName, len(params)+1))
		params = append(params, value)
	}

	if data.State != nil {
		addParam("state", data.State)
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
