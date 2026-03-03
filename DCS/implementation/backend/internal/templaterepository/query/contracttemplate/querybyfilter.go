package contracttemplate

import (
	"context"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/base/event"
	templaterepository2 "digital-contracting-service/internal/templaterepository/datatype/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/datatype/templatetype"
	"digital-contracting-service/internal/templaterepository/db"
	templateevents "digital-contracting-service/internal/templaterepository/event"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type GetAllMetadataByFilterQry struct {
	RetrievedBy string

	DID            *string
	DocumentNumber *int
	Version        *int
	State          *templatestate.TemplateState
	TemplateType   *templatetype.TemplateType
	Name           *string
	Description    *string
	Filter         *string
}

type GetAllMetadataByFilterResult struct {
	DID            string
	DocumentNumber int
	Version        int
	State          templatestate.TemplateState
	TemplateType   templatetype.TemplateType
	Name           string
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	MetaData       datatype.JSON
}

type GetAllMetaDataByFilterHandler struct {
	Ctx    context.Context
	DB     *sqlx.DB
	CTRepo db.TemplateRepository
}

func (h *GetAllMetaDataByFilterHandler) Handle(query GetAllMetadataByFilterQry) ([]GetAllMetadataByFilterResult, error) {

	ctx, cancel := context.WithTimeout(h.Ctx, base.TransactionTimeout())
	defer cancel()

	tx, err := h.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create transaction: %w", err)
	}
	defer tx.Rollback()

	searchValues := templaterepository2.SearchValues{
		DID:            query.DID,
		DocumentNumber: query.DocumentNumber,
		Version:        query.Version,
		State:          query.State,
		TemplateType:   query.TemplateType,
		Name:           query.Name,
		Description:    query.Description,
		Filter:         query.Filter,
	}

	contractTemplates, err := h.CTRepo.ReadAllMetaDataByFilter(tx, searchValues)
	if err != nil {
		return nil, fmt.Errorf("could not read all contract templates: %w", err)
	}

	evt := templateevents.RetrieveAllEvent{
		RetrievedBy: query.RetrievedBy,
		OccurredAt:  time.Now(),
	}
	err = event.Create(h.Ctx, tx, evt)
	if err != nil {
		return nil, fmt.Errorf("could not create event: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	result := make([]GetAllMetadataByFilterResult, len(contractTemplates))
	for i, data := range contractTemplates {
		result[i] = GetAllMetadataByFilterResult{
			DID:            data.DID,
			DocumentNumber: data.DocumentNumber,
			Version:        data.Version,
			State:          data.State,
			TemplateType:   data.TemplateType,
			Name:           *data.Name,
			Description:    *data.Description,
			CreatedAt:      data.CreatedAt,
			UpdatedAt:      data.UpdatedAt,
		}
	}

	return result, nil
}
