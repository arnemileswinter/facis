package templaterepository

import (
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/datatype/templatetype"
	"time"
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

type ProcessData struct {
	DID            string                      `db:"did"`
	DocumentNumber int                         `db:"document_number"`
	Version        int                         `db:"version"`
	State          templatestate.TemplateState `db:"state"`
	CreatedBy      string                      `db:"created_by"`
	UpdatedAt      time.Time                   `db:"updated_at"`
}

type UpdateData struct {
	DID            string                       `db:"did"`
	DocumentNumber int                          `db:"document_number"`
	Version        int                          `db:"version"`
	State          *templatestate.TemplateState `db:"state"`
	TemplateType   *templatetype.TemplateType   `db:"template_type"`
	Name           *string                      `db:"name"`
	Description    *string                      `db:"description"`
	TemplateData   *datatype.JSON               `db:"template_data"`
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
