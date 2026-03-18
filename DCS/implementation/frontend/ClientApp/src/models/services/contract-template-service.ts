import type { ContractTemplate } from '../contract-template'
import type {
  ContractTemplateCreateRequest,
  ContractTemplateSubmitRequest,
  ContractTemplateUpdateRequest,
  ContractTemplateSearchRequest,
  ContractTemplateRetrieveRequest,
  ContractTemplateRetrieveByIdRequest,
  ContractTemplateApproveRequest,
  ContractTemplateRejectRequest,
  ContractTemplateVerifyRequest,
} from '../requests/template-request'
import type {
  ContractTemplateCreateResponse,
  ContractTemplateSubmitResponse,
  ContractTemplateUpdateResponse,
  ContractTemplateSearchResponse,
  ContractTemplateRetrieveResponse,
  ContractTemplateApproveResponse,
  ContractTemplateRejectResponse,
  ContractTemplateVerifyResponse,
} from '../responses/template-response'

export interface ContractTemplateService {
  create: (request: ContractTemplateCreateRequest) => Promise<ContractTemplateCreateResponse>
  submit: (request: ContractTemplateSubmitRequest) => Promise<ContractTemplateSubmitResponse>
  update: (request: ContractTemplateUpdateRequest) => Promise<ContractTemplateUpdateResponse>
  search: (request: ContractTemplateSearchRequest) => Promise<ContractTemplateSearchResponse>
  retrieve: (request?: ContractTemplateRetrieveRequest) => Promise<ContractTemplateRetrieveResponse>
  retrieveById: (request: ContractTemplateRetrieveByIdRequest) => Promise<ContractTemplate | null>
  approve: (request: ContractTemplateApproveRequest) => Promise<ContractTemplateApproveResponse>
  reject: (request: ContractTemplateRejectRequest) => Promise<ContractTemplateRejectResponse>
  verify: (request: ContractTemplateVerifyRequest) => Promise<ContractTemplateVerifyResponse>
}
