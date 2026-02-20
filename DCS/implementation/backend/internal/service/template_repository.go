package service

import (
	"context"
	templaterepository "digital-contracting-service/gen/template_repository"
	"digital-contracting-service/internal/auth"

	"goa.design/clue/log"
)

type templateRepositorysrvc struct {
	auth.JWTAuthenticator
}

func NewTemplateRepository(jwtAuth auth.JWTAuthenticator) templaterepository.Service {
	return &templateRepositorysrvc{JWTAuthenticator: jwtAuth}
}

func (s *templateRepositorysrvc) Create(ctx context.Context, p *templaterepository.CreatePayload) (res string, err error) {
	log.Printf(ctx, "templateRepository.create")
	return
}

func (s *templateRepositorysrvc) Submit(ctx context.Context, p *templaterepository.SubmitPayload) (res string, err error) {
	log.Printf(ctx, "templateRepository.submit")
	return
}

func (s *templateRepositorysrvc) Update(ctx context.Context, p *templaterepository.UpdatePayload) (res int, err error) {
	log.Printf(ctx, "templateRepository.update")
	return
}

func (s *templateRepositorysrvc) UpdateManage(ctx context.Context, p *templaterepository.UpdateManagePayload) (res int, err error) {
	log.Printf(ctx, "templateRepository.update_manage")
	return
}

func (s *templateRepositorysrvc) Search(ctx context.Context, p *templaterepository.SearchPayload) (res []any, err error) {
	log.Printf(ctx, "templateRepository.search")
	return
}

func (s *templateRepositorysrvc) Retrieve(ctx context.Context, p *templaterepository.RetrievePayload) (res any, err error) {
	log.Printf(ctx, "templateRepository.retrieve")
	return
}

func (s *templateRepositorysrvc) RetrieveByID(ctx context.Context, p *templaterepository.RetrieveByIDPayload) (res any, err error) {
	log.Printf(ctx, "templateRepository.retrieve_by_id")
	return
}

func (s *templateRepositorysrvc) Verify(ctx context.Context, p *templaterepository.VerifyPayload) (res any, err error) {
	log.Printf(ctx, "templateRepository.verify")
	return
}

func (s *templateRepositorysrvc) Approve(ctx context.Context, p *templaterepository.ApprovePayload) (res int, err error) {
	log.Printf(ctx, "templateRepository.approve")
	return
}

func (s *templateRepositorysrvc) Reject(ctx context.Context, p *templaterepository.RejectPayload) (res int, err error) {
	log.Printf(ctx, "templateRepository.reject")
	return
}

func (s *templateRepositorysrvc) Register(ctx context.Context, p *templaterepository.RegisterPayload) (res any, err error) {
	log.Printf(ctx, "templateRepository.register")
	return
}

func (s *templateRepositorysrvc) Archive(ctx context.Context, p *templaterepository.ArchivePayload) (res int, err error) {
	log.Printf(ctx, "templateRepository.archive")
	return
}

func (s *templateRepositorysrvc) Audit(ctx context.Context, p *templaterepository.AuditPayload) (res []string, err error) {
	log.Printf(ctx, "templateRepository.audit")
	return
}
