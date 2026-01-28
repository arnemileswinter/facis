package service

import (
	"context"
	templaterepository "digital-contracting-service/gen/template_repository"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/command"
	"digital-contracting-service/internal/template_repository/query"

	"github.com/jmoiron/sqlx"
	"goa.design/clue/log"
)

// TemplateRepository service example implementation.
// The example methods log the requests and return zero values.
type templateRepositorysrvc struct {
	db *sqlx.DB
}

// NewTemplateRepository returns the TemplateRepository service implementation.
func NewTemplateRepository(ctx context.Context, db *sqlx.DB) (templaterepository.Service, error) {
	return &templateRepositorysrvc{
		db: db,
	}, nil
}

// Create a new template.
func (s *templateRepositorysrvc) Create(ctx context.Context, req *templaterepository.ContractTemplateCreateRequest) (*templaterepository.ContractTemplateCreateResponse, error) {

	jsonMetaData, err := datatype.NewJSON(req.MetaData)
	if err != nil {
		log.Errorf(ctx, err, "failed to convert metadata")
		return nil, templaterepository.MakeInternalError(err)
	}

	did, err := base.GetDID()
	if err != nil {
		log.Errorf(ctx, err, "failed to generate uuid")
		return nil, templaterepository.MakeInternalError(err)
	}

	cmd := command.CreateTemplateContractCommand{
		DID:         *did,
		CreatedBy:   "",
		Name:        req.Name,
		Description: req.Description,
		MetaData:    &jsonMetaData,
	}
	createHandler := command.CreateTemplateContractHandler{
		Db: s.db,
	}
	err = createHandler.Handle(cmd)
	if err != nil {
		log.Errorf(ctx, err, "failed to create template contract")
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateCreateResponse{
		Did: *did,
	}, nil
}

// with action flag { forwardTo: "approval" | "draft" } and optional
// reviewComments. allow resubmission path with approver comments.
func (s *templateRepositorysrvc) Submit(ctx context.Context) (res string, err error) {
	log.Printf(ctx, "templateRepository.submit")
	return
}

// persist reviewer edits (metadata/clauses/semantics).
func (s *templateRepositorysrvc) Update(ctx context.Context, req *templaterepository.ContractTemplateUpdateRequest) (res *templaterepository.ContractTemplateUpdateResponse, err error) {

	metaData, err := datatype.NewJSON(req.MetaData)
	if err != nil {
		log.Errorf(ctx, err, "failed to convert metadata")
		return nil, templaterepository.MakeInternalError(err)
	}
	cmd := command.UpdateTemplateContractCommand{
		DID:         req.Did,
		UpdatedBy:   "",
		Name:        req.Name,
		Description: req.Description,
		MetaData:    &metaData,
	}
	handler := command.UpdateTemplateContractHandler{
		Db: s.db,
	}
	err = handler.Handle(cmd)
	if err != nil {
		log.Errorf(ctx, err, "failed to update template contract")
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateUpdateResponse{
		Did: req.Did,
	}, nil
}

// update metadata or status.
func (s *templateRepositorysrvc) UpdateManage(ctx context.Context) (res int, err error) {
	log.Printf(ctx, "templateRepository.update_manage")
	return
}

// perform filtered searches.
func (s *templateRepositorysrvc) Search(ctx context.Context) (res []any, err error) {
	log.Printf(ctx, "templateRepository.search")
	return
}

// load submitted template and history/provenance summary. fetch reviewed
// template with metadata, review history, and validation results. fetch all
// template entries for dashboard view.
func (s *templateRepositorysrvc) Retrieve(ctx context.Context) (res any, err error) {
	log.Printf(ctx, "templateRepository.retrieve")
	return
}

// Retrieve a template by template id.
func (s *templateRepositorysrvc) RetrieveByID(ctx context.Context, req *templaterepository.ContractTemplateRetrieveByIDRequest) (res *templaterepository.ContractTemplateRetrieveByIDResponse, err error) {

	qry := query.RetrieveContractByIdQuery{
		DID:         req.TemplateID,
		RetrievedBy: "",
	}
	queryHandler := query.RetrieveContractTemplateByIdHandler{
		Ctx: ctx,
		Db:  s.db,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		log.Errorf(ctx, err, "failed to get template contract")
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateRetrieveByIDResponse{
		Did:            contractTemplate.DID,
		DocumentNumber: contractTemplate.DocumentNumber,
		Version:        contractTemplate.Version,
		State:          contractTemplate.State.String(),
		Name:           &contractTemplate.Name,
		Description:    &contractTemplate.Description,
		CreatedBy:      contractTemplate.CreatedBy,
		CreatedAt:      contractTemplate.CreatedAt.String(),
		MetaData:       contractTemplate.MetaData,
	}, nil
}

// run policy, schema, and semantic validations; return findings.
func (s *templateRepositorysrvc) Verify(ctx context.Context) (res any, err error) {
	log.Printf(ctx, "templateRepository.verify")
	return
}

// mark template as approved, with optional decision notes.
func (s *templateRepositorysrvc) Approve(ctx context.Context) (res int, err error) {
	log.Printf(ctx, "templateRepository.approve")
	return
}

// mark template as rejected, requiring reason field.
func (s *templateRepositorysrvc) Reject(ctx context.Context) (res int, err error) {
	log.Printf(ctx, "templateRepository.reject")
	return
}

// register new template into the repository.
func (s *templateRepositorysrvc) Register(ctx context.Context) (res any, err error) {
	log.Printf(ctx, "templateRepository.register")
	return
}

// archive obsolete template.
func (s *templateRepositorysrvc) Archive(ctx context.Context) (res int, err error) {
	log.Printf(ctx, "templateRepository.archive")
	return
}

// retrieve audit history of template actions.
func (s *templateRepositorysrvc) Audit(ctx context.Context) (res []string, err error) {
	log.Printf(ctx, "templateRepository.audit")
	return
}
