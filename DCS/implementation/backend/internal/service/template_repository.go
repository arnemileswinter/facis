package service

import (
	"context"
	templaterepository "digital-contracting-service/gen/template_repository"
	"digital-contracting-service/internal/auth"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/middleware"
	"digital-contracting-service/internal/templaterepository/command"
	"digital-contracting-service/internal/templaterepository/datatype/actionflag"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"digital-contracting-service/internal/templaterepository/datatype/templatetype"
	"digital-contracting-service/internal/templaterepository/query/contracttemplate"
	"time"

	"github.com/jmoiron/sqlx"
	"goa.design/clue/log"
)

// TemplateRepository service example implementation.
// The example methods log the requests and return zero values.
type templateRepositorysrvc struct {
	DB *sqlx.DB
	auth.JWTAuthenticator
}

// NewTemplateRepository returns the TemplateRepository service implementation.
func NewTemplateRepository(ctx context.Context, db *sqlx.DB, jwtAuth auth.JWTAuthenticator) (templaterepository.Service, error) {
	return &templateRepositorysrvc{
		DB:               db,
		JWTAuthenticator: jwtAuth,
	}, nil
}

// Create a new template.
func (s *templateRepositorysrvc) Create(ctx context.Context, req *templaterepository.ContractTemplateCreateRequest) (*templaterepository.ContractTemplateCreateResponse, error) {

	templateType, err := templatetype.NewTemplateType(req.TemplateType)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	jsonMetaData, err := datatype.NewJSON(req.TemplateData)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	did, err := base.GetDID()
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	cmd := command.CreateCommand{
		DID:          *did,
		CreatedBy:    middleware.GetUsername(ctx),
		TemplateType: templateType,
		Name:         req.Name,
		Description:  req.Description,
		TemplateData: &jsonMetaData,
	}
	createHandler := command.CreateHandler{
		Ctx: ctx,
		DB:  s.DB,
	}
	err = createHandler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateCreateResponse{
		Did:            *did,
		DocumentNumber: 1,
		Version:        1,
	}, nil
}

// with action flag { forwardTo: "approval" | "draft" } and optional
// reviewComments. allow resubmission path with approver comments.
func (s *templateRepositorysrvc) Submit(ctx context.Context, req *templaterepository.ContractTemplateSubmitRequest) (res *templaterepository.ContractTemplateSubmitResponse, err error) {

	updatedAt, err := time.Parse(time.RFC3339, req.UpdatedAt)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	var actionFlag *actionflag.ActionFlag
	if req.ForwardTo != nil {
		flag, err := actionflag.NewActionFlag(*req.ForwardTo)
		if err != nil {
			return nil, templaterepository.MakeInternalError(err)
		}
		actionFlag = &flag
	}

	cmd := command.SubmitCommand{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		UpdatedAt:      updatedAt,
		SubmittedBy:    middleware.GetUsername(ctx),
		ActionFlag:     actionFlag,
		Comments:       req.Comments,
	}
	handler := command.SubmitHandler{
		DB: s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateSubmitResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

// persist reviewer edits (metadata/clauses/semantics).
func (s *templateRepositorysrvc) Update(ctx context.Context, req *templaterepository.ContractTemplateUpdateRequest) (res *templaterepository.ContractTemplateUpdateResponse, err error) {

	updatedAt, err := time.Parse(time.RFC3339, req.UpdatedAt)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	metaData, err := datatype.NewJSON(req.TemplateData)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	var templateType *templatetype.TemplateType
	if req.TemplateType != nil {
		tType, err := templatetype.NewTemplateType(*req.TemplateType)
		if err != nil {
			return nil, templaterepository.MakeInternalError(err)
		}
		templateType = &tType
	}

	cmd := command.UpdateCommand{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		UpdatedAt:      updatedAt,
		TemplateType:   templateType,
		Name:           req.Name,
		Description:    req.Description,
		TemplateData:   &metaData,
	}
	handler := command.UpdateHandler{
		Ctx: ctx,
		DB:  s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateUpdateResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

func (s *templateRepositorysrvc) UpdateManage(ctx context.Context, p *templaterepository.UpdateManagePayload) (res int, err error) {
	log.Printf(ctx, "templateRepository.update_manage")
	return
}

// perform filtered searches.
func (s *templateRepositorysrvc) Search(ctx context.Context, req *templaterepository.ContractTemplateSearchRequest) (res []*templaterepository.ContractTemplateSearchResponse, err error) {

	var state *templatestate.TemplateState
	if req.State != nil {
		tState, err := templatestate.NewTemplateState(*req.State)
		if err != nil {
			return nil, templaterepository.MakeInternalError(err)
		}

		state = &tState
	}

	qry := contracttemplate.GetAllMetaDataByFilterQuery{
		RetrievedBy:    middleware.GetUsername(ctx),
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		State:          state,
		Name:           req.Name,
		Description:    req.Description,
		Filter:         req.Filter,
	}
	queryHandler := contracttemplate.GetAllMetaDataByFilterHandler{
		Ctx: ctx,
		DB:  s.DB,
	}
	result, err := queryHandler.Handle(qry)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	var contractTemplates []*templaterepository.ContractTemplateSearchResponse
	for _, item := range result {
		contractTemplates = append(contractTemplates, &templaterepository.ContractTemplateSearchResponse{
			Did:            item.DID,
			DocumentNumber: item.DocumentNumber,
			Version:        item.Version,
			State:          item.State.String(),
			Name:           &item.Name,
			Description:    &item.Description,
			CreatedAt:      item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      item.UpdatedAt.Format(time.RFC3339),
		})
	}

	return contractTemplates, nil
}

// retrieve templates
func (s *templateRepositorysrvc) Retrieve(ctx context.Context, req *templaterepository.ContractTemplateRetrieveRequest) (res *templaterepository.ContractTemplateRetrieveResponse, err error) {

	qry := contracttemplate.GetAllMetaDataQuery{
		RetrievedBy: middleware.GetUsername(ctx),
	}
	queryHandler := contracttemplate.GetAllMetaDataHandler{
		Ctx: ctx,
		DB:  s.DB,
	}
	result, err := queryHandler.Handle(qry)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	var contractTemplates []*templaterepository.ContractTemplateItem
	for _, item := range result.ContractTemplates {
		contractTemplates = append(contractTemplates, &templaterepository.ContractTemplateItem{
			Did:            item.DID,
			DocumentNumber: item.DocumentNumber,
			Version:        item.Version,
			State:          item.State.String(),
			Name:           &item.Name,
			Description:    &item.Description,
			CreatedAt:      item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      item.UpdatedAt.Format(time.RFC3339),
		})
	}

	var reviewTasks []*templaterepository.ContractTemplateReviewTaskItem
	for _, item := range result.ReviewerTasks {
		reviewTasks = append(reviewTasks, &templaterepository.ContractTemplateReviewTaskItem{
			Did:            item.DID,
			DocumentNumber: item.DocumentNumber,
			Version:        item.Version,
			Reviewer:       item.Reviewer,
			State:          item.State.String(),
			CreatedAt:      item.CreatedAt.Format(time.RFC3339),
		})
	}

	var approvalTasks []*templaterepository.ContractTemplateApprovalTaskItem
	for _, item := range result.ApprovalTasks {
		approvalTasks = append(approvalTasks, &templaterepository.ContractTemplateApprovalTaskItem{
			Did:            item.DID,
			DocumentNumber: item.DocumentNumber,
			Version:        item.Version,
			State:          item.State.String(),
			Approver:       item.Approver,
			CreatedAt:      item.CreatedAt.Format(time.RFC3339),
		})
	}

	return &templaterepository.ContractTemplateRetrieveResponse{
		ContractTemplates: contractTemplates,
		ReviewTasks:       reviewTasks,
		ApprovalTasks:     approvalTasks,
	}, nil
}

// Retrieve a template by template id.
func (s *templateRepositorysrvc) RetrieveByID(ctx context.Context, req *templaterepository.ContractTemplateRetrieveByIDRequest) (res *templaterepository.ContractTemplateRetrieveByIDResponse, err error) {

	qry := contracttemplate.GetByIDQuery{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		RetrievedBy:    middleware.GetUsername(ctx),
	}
	queryHandler := contracttemplate.GetByIDHandler{
		Ctx: ctx,
		DB:  s.DB,
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
		Name:           contractTemplate.Name,
		Description:    contractTemplate.Description,
		CreatedBy:      contractTemplate.CreatedBy,
		CreatedAt:      contractTemplate.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      contractTemplate.UpdatedAt.Format(time.RFC3339),
		TemplateData:   contractTemplate.TemplateData,
	}, nil
}

// run policy, schema, and semantic validations; return findings.
func (s *templateRepositorysrvc) Verify(ctx context.Context, req *templaterepository.ContractTemplateVerifyRequest) (res *templaterepository.ContractTemplateVerifyResponse, err error) {

	updatedAt, err := time.Parse(time.RFC3339, req.UpdatedAt)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	cmd := command.VerifyCommand{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		UpdatedAt:      updatedAt,
	}
	handler := command.VerifyHandler{
		DB: s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateVerifyResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

// mark template as approved, with optional decision notes.
func (s *templateRepositorysrvc) Approve(ctx context.Context, req *templaterepository.ContractTemplateApproveRequest) (res *templaterepository.ContractTemplateApproveResponse, err error) {

	updatedAt, err := time.Parse(time.RFC3339, req.UpdatedAt)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	cmd := command.ApproveCommand{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		UpdatedAt:      updatedAt,
		ApprovedBy:     middleware.GetUsername(ctx),
		DecisionNotes:  req.DecisionNotes,
	}
	handler := command.ApproveHandler{
		DB: s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateApproveResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

// mark template as rejected, requiring reason field.
func (s *templateRepositorysrvc) Reject(ctx context.Context, req *templaterepository.ContractTemplateRejectRequest) (res *templaterepository.ContractTemplateRejectResponse, err error) {

	updatedAt, err := time.Parse(time.RFC3339, req.UpdatedAt)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	cmd := command.RejectCommand{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		UpdatedAt:      updatedAt,
		RejectedBy:     middleware.GetUsername(ctx),
		Reason:         req.Reason,
	}
	handler := command.RejectHandler{
		Ctx: ctx,
		DB:  s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateRejectResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

// register new template into the repository.
func (s *templateRepositorysrvc) Register(ctx context.Context, req *templaterepository.ContractTemplateRegisterRequest) (res *templaterepository.ContractTemplateRegisterResponse, err error) {

	updatedAt, err := time.Parse(time.RFC3339, req.UpdatedAt)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	cmd := command.RegisterCommand{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		UpdatedAt:      updatedAt,
		RegisteredBy:   middleware.GetUsername(ctx),
	}
	handler := command.RegisterHandler{
		Ctx: ctx,
		DB:  s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateRegisterResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

// archive obsolete template.
func (s *templateRepositorysrvc) Archive(ctx context.Context, req *templaterepository.ContractTemplateArchiveRequest) (res *templaterepository.ContractTemplateArchiveResponse, err error) {

	updatedAt, err := time.Parse(time.RFC3339, req.UpdatedAt)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	cmd := command.ArchiveCommand{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		UpdatedAt:      updatedAt,
		ArchivedBy:     middleware.GetUsername(ctx),
	}
	handler := command.ArchiveHandler{
		Ctx: ctx,
		DB:  s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ContractTemplateArchiveResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

func (s *templateRepositorysrvc) Audit(ctx context.Context, p *templaterepository.AuditPayload) (res []string, err error) {
	log.Printf(ctx, "templateRepository.audit")
	return
}
