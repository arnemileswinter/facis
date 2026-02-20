import type { ActionFlag } from '../../types/action-flag'

interface ContractTemplateBaseRequest {
  did: string
  document_number: number
  version: number
}

export interface ContractTemplateCreateRequest {
  name?: string
  description?: string
  /** The template data of the contract template */
  template_data?: any
}

export interface ContractTemplateSubmitRequest extends ContractTemplateBaseRequest {
  updated_at: string
  forward_to?: ActionFlag
  comments?: string[]
}

export interface ContractTemplateUpdateRequest extends ContractTemplateBaseRequest {
  updated_at: string
  name?: string
  description?: string
  /** The template data of the contract template */
  template_data?: any
}

export interface ContractTemplateSearchRequest {
  filter: any
}

export interface ContractTemplateRetrieveRequest {}

export interface ContractTemplateRetrieveByIdRequest extends ContractTemplateBaseRequest {}

export interface ContractTemplateApproveRequest extends ContractTemplateBaseRequest {
  updated_at: string
  decision_notes?: string[]
}

export interface ContractTemplateRejectRequest extends ContractTemplateBaseRequest {
  updated_at: string
  /** Reason for rejecting the contract template */
  reason: string
}
