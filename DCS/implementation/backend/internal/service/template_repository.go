package service

import (
	"context"
	templaterepository "digital-contracting-service/gen/template_repository"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/command"
	"digital-contracting-service/internal/template_repository/datatype/action_flag"
	"digital-contracting-service/internal/template_repository/query"

	"github.com/jmoiron/sqlx"
	"goa.design/clue/log"
)

// TemplateRepository service example implementation.
// The example methods log the requests and return zero values.
type templateRepositorysrvc struct {
	DB *sqlx.DB
}

// NewTemplateRepository returns the TemplateRepository service implementation.
func NewTemplateRepository(ctx context.Context, db *sqlx.DB) (templaterepository.Service, error) {
	return &templateRepositorysrvc{
		DB: db,
	}, nil
}

// Create a new template.
func (s *templateRepositorysrvc) Create(ctx context.Context, req *templaterepository.ContractTemplateCreateRequest) (*templaterepository.ContractTemplateCreateResponse, error) {

	jsonMetaData, err := datatype.NewJSON(req.MetaData)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	did, err := base.GetDID()
	if err != nil {
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
		Ctx: ctx,
		DB:  s.DB,
	}
	err = createHandler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateCreateResponse{
		Did: *did,
	}, nil
}

// with action flag { forwardTo: "approval" | "draft" } and optional
// reviewComments. allow resubmission path with approver comments.
func (s *templateRepositorysrvc) Submit(ctx context.Context, req *templaterepository.ContractTemplateSubmitRequest) (res *templaterepository.ContractTemplateSubmitResponse, err error) {

	var actionFlag *action_flag.ActionFlag
	if req.ForwardTo != nil {
		flag, err := action_flag.NewActionFlag(*req.ForwardTo)
		if err != nil {
			return nil, templaterepository.MakeInternalError(err)
		}
		actionFlag = &flag
	}

	cmd := command.SubmitTemplateContractCommand{
		DID:            req.Did,
		SubmittedBy:    "",
		ActionFlag:     actionFlag,
		ReviewComments: req.ReviewComments,
	}
	handler := command.SubmitTemplateContractHandler{
		DB: s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateSubmitResponse{
		Did: req.Did,
	}, nil
}

// persist reviewer edits (metadata/clauses/semantics).
func (s *templateRepositorysrvc) Update(ctx context.Context, req *templaterepository.ContractTemplateUpdateRequest) (res *templaterepository.ContractTemplateUpdateResponse, err error) {

	metaData, err := datatype.NewJSON(req.MetaData)
	if err != nil {
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
		Ctx: ctx,
		DB:  s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
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
func (s *templateRepositorysrvc) Retrieve(ctx context.Context) (res []*templaterepository.ContractTemplateRetrieveResponse, err error) {

	qry := query.GetAllContractTemplatesQuery{
		RetrievedBy: "",
	}
	queryHandler := query.GetAllContractTemplateHandler{
		Ctx: ctx,
		Db:  s.DB,
	}
	result, err := queryHandler.Handle(qry)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	var contractTemplates []*templaterepository.ContractTemplateRetrieveResponse
	for _, item := range result {
		contractTemplates = append(contractTemplates, &templaterepository.ContractTemplateRetrieveResponse{
			Did:            item.DID,
			DocumentNumber: item.DocumentNumber,
			Version:        item.Version,
			State:          item.State.String(),
			Name:           &item.Name,
			Description:    &item.Description,
			CreatedBy:      item.CreatedBy,
			CreatedAt:      item.CreatedAt.String(),
			MetaData:       item.MetaData,
		})
	}

	return contractTemplates, nil
}

// Retrieve a template by template id.
func (s *templateRepositorysrvc) RetrieveByID(ctx context.Context, req *templaterepository.ContractTemplateRetrieveByIDRequest) (res *templaterepository.ContractTemplateRetrieveByIDResponse, err error) {

	qry := query.GetContractTemplateByIdQuery{
		DID:         req.TemplateID,
		RetrievedBy: "",
	}
	queryHandler := query.GetContractTemplateByIdHandler{
		Ctx: ctx,
		Db:  s.DB,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
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
func (s *templateRepositorysrvc) Approve(ctx context.Context, req *templaterepository.ContractTemplateApproveRequest) (res *templaterepository.ContractTemplateApproveResponse, err error) {

	cmd := command.ApproveTemplateContractCommand{
		DID:           req.Did,
		ApprovedBy:    "",
		DecisionNotes: req.DecisionNotes,
	}
	handler := command.ApproveTemplateContractHandler{
		DB: s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateApproveResponse{
		Did: req.Did,
	}, nil

}

// mark template as rejected, requiring reason field.
func (s *templateRepositorysrvc) Reject(ctx context.Context, req *templaterepository.ContractTemplateRejectRequest) (res *templaterepository.ContractTemplateRejectResponse, err error) {

	cmd := command.RejectTemplateContractCommand{
		DID:        req.Did,
		RejectedBy: "",
		Reason:     req.Reason,
	}
	handler := command.RejectTemplateContractHandler{
		Ctx: ctx,
		DB:  s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateRejectResponse{
		Did: req.Did,
	}, nil

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
