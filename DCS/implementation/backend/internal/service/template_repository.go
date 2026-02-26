package service

import (
	"context"
	templaterepository "digital-contracting-service/gen/template_repository"
	"digital-contracting-service/internal/auth"
	"digital-contracting-service/internal/base"
	"digital-contracting-service/internal/base/datatype"
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
		DB: db,
		JWTAuthenticator: jwtAuth,
	}, nil
}

// Create a new template.
func (s *templateRepositorysrvc) Create(ctx context.Context, req *templaterepository.CreateRequest) (*templaterepository.CreateResponse, error) {

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
		CreatedBy:    "Test User",
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

	return &templaterepository.CreateResponse{
		Did:            *did,
		DocumentNumber: 1,
		Version:        1,
	}, nil
}

// with action flag { forwardTo: "approval" | "draft" } and optional
// reviewComments. allow resubmission path with approver comments.
func (s *templateRepositorysrvc) Submit(ctx context.Context, req *templaterepository.SubmitRequest) (res *templaterepository.SubmitResponse, err error) {

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
		SubmittedBy:    "Test User",
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

	return &templaterepository.SubmitResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

// persist reviewer edits (metadata/clauses/semantics).
func (s *templateRepositorysrvc) Update(ctx context.Context, req *templaterepository.UpdateRequest) (res *templaterepository.UpdateResponse, err error) {

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

	return &templaterepository.UpdateResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

// update metadata or status.
func (s *templateRepositorysrvc) UpdateManage(ctx context.Context, req *templaterepository.UpdateManageRequest) (res *templaterepository.UpdateManageResponse, err error) {

	updatedAt, err := time.Parse(time.RFC3339, req.UpdatedAt)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	metaData, err := datatype.NewJSON(req.TemplateData)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	var state *templatestate.TemplateState
	if req.State != nil {
		ts, err := templatestate.NewTemplateState(*req.State)
		if err != nil {
			return nil, templaterepository.MakeInternalError(err)
		}
		state = &ts
	}

	var templateType *templatetype.TemplateType
	if req.TemplateType != nil {
		tType, err := templatetype.NewTemplateType(*req.TemplateType)
		if err != nil {
			return nil, templaterepository.MakeInternalError(err)
		}
		templateType = &tType
	}

	cmd := command.UpdateManageCommand{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		State:          state,
		UpdatedAt:      updatedAt,
		TemplateType:   templateType,
		Name:           req.Name,
		Description:    req.Description,
		TemplateData:   &metaData,
	}
	handler := command.UpdateManageHandler{
		Ctx: ctx,
		DB:  s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.UpdateManageResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

// perform filtered searches.
func (s *templateRepositorysrvc) Search(ctx context.Context, req *templaterepository.SearchRequest) (res []*templaterepository.SearchResponse, err error) {

	var state *templatestate.TemplateState
	if req.State != nil {
		tState, err := templatestate.NewTemplateState(*req.State)
		if err != nil {
			return nil, templaterepository.MakeInternalError(err)
		}

		state = &tState
	}

	qry := contracttemplate.GetAllMetaDataByFilterQuery{
		RetrievedBy:    "Test User",
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

	var contractTemplates []*templaterepository.SearchResponse
	for _, item := range result {
		contractTemplates = append(contractTemplates, &templaterepository.SearchResponse{
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
func (s *templateRepositorysrvc) Retrieve(ctx context.Context, req *templaterepository.RetrieveRequest) (res *templaterepository.RetrieveResponse, err error) {

	qry := contracttemplate.GetAllMetaDataQuery{
		RetrievedBy: "Test User",
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

	var reviewTasks []*templaterepository.ReviewTaskItem
	for _, item := range result.ReviewerTasks {
		reviewTasks = append(reviewTasks, &templaterepository.ReviewTaskItem{
			Did:            item.DID,
			DocumentNumber: item.DocumentNumber,
			Version:        item.Version,
			Reviewer:       item.Reviewer,
			State:          item.State.String(),
			CreatedAt:      item.CreatedAt.Format(time.RFC3339),
		})
	}

	var approvalTasks []*templaterepository.ApprovalTaskItem
	for _, item := range result.ApprovalTasks {
		approvalTasks = append(approvalTasks, &templaterepository.ApprovalTaskItem{
			Did:            item.DID,
			DocumentNumber: item.DocumentNumber,
			Version:        item.Version,
			State:          item.State.String(),
			Approver:       item.Approver,
			CreatedAt:      item.CreatedAt.Format(time.RFC3339),
		})
	}

	return &templaterepository.RetrieveResponse{
		ContractTemplates: contractTemplates,
		ReviewTasks:       reviewTasks,
		ApprovalTasks:     approvalTasks,
	}, nil
}

// Retrieve a template by template id.
func (s *templateRepositorysrvc) RetrieveByID(ctx context.Context, req *templaterepository.RetrieveByIDRequest) (res *templaterepository.RetrieveByIDResponse, err error) {

	qry := contracttemplate.GetByIDQuery{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		RetrievedBy:    "Test User",
	}
	queryHandler := contracttemplate.GetByIDHandler{
		Ctx: ctx,
		DB:  s.DB,
	}
	contractTemplate, err := queryHandler.Handle(qry)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.RetrieveByIDResponse{
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
func (s *templateRepositorysrvc) Verify(ctx context.Context, req *templaterepository.VerifyRequest) (res *templaterepository.VerifyResponse, err error) {

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

	return &templaterepository.VerifyResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

// mark template as approved, with optional decision notes.
func (s *templateRepositorysrvc) Approve(ctx context.Context, req *templaterepository.ApproveRequest) (res *templaterepository.ApproveResponse, err error) {

	updatedAt, err := time.Parse(time.RFC3339, req.UpdatedAt)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	cmd := command.ApproveCommand{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		UpdatedAt:      updatedAt,
		ApprovedBy:     "Test User",
		DecisionNotes:  req.DecisionNotes,
	}
	handler := command.ApproveHandler{
		DB: s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ApproveResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

// mark template as rejected, requiring reason field.
func (s *templateRepositorysrvc) Reject(ctx context.Context, req *templaterepository.RejectRequest) (res *templaterepository.RejectResponse, err error) {

	updatedAt, err := time.Parse(time.RFC3339, req.UpdatedAt)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	cmd := command.RejectCommand{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		UpdatedAt:      updatedAt,
		RejectedBy:     "Test User",
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

	return &templaterepository.RejectResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

// register new template into the repository.
func (s *templateRepositorysrvc) Register(ctx context.Context, req *templaterepository.RegisterRequest) (res *templaterepository.RegisterResponse, err error) {

	updatedAt, err := time.Parse(time.RFC3339, req.UpdatedAt)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	cmd := command.RegisterCommand{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		UpdatedAt:      updatedAt,
		RegisteredBy:   "Test User",
	}
	handler := command.RegisterHandler{
		Ctx: ctx,
		DB:  s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.RegisterResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

// archive obsolete template.
func (s *templateRepositorysrvc) Archive(ctx context.Context, req *templaterepository.ArchiveRequest) (res *templaterepository.ArchiveResponse, err error) {

	updatedAt, err := time.Parse(time.RFC3339, req.UpdatedAt)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	cmd := command.ArchiveCommand{
		DID:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
		UpdatedAt:      updatedAt,
		ArchivedBy:     "Test User",
	}
	handler := command.ArchiveHandler{
		Ctx: ctx,
		DB:  s.DB,
	}
	err = handler.Handle(cmd)
	if err != nil {
		return nil, templaterepository.MakeInternalError(err)
	}

	return &templaterepository.ArchiveResponse{
		Did:            req.Did,
		DocumentNumber: req.DocumentNumber,
		Version:        req.Version,
	}, nil
}

// retrieve audit history of template actions.
func (s *templateRepositorysrvc) Audit(ctx context.Context) (res []string, err error) {
	log.Printf(ctx, "templateRepository.audit")
	return
}
