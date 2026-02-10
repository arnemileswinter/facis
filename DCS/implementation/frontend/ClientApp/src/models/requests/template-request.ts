import type { ActionFlag } from '../../types/action-flag'

interface ContractTemplateBaseRequest {
  did: string
}

export interface ContractTemplateCreateRequest {
  name?: string
  description?: string
  /** The metadata of the contract template */
  meta_data?: any
}

export interface ContractTemplateSubmitRequest extends ContractTemplateBaseRequest {
  forward_to?: ActionFlag
  review_comments?: string[]
}

export interface ContractTemplateUpdateRequest extends ContractTemplateBaseRequest {
  name?: string
  description?: string
  /** The metadata of the contract template */
  meta_data?: any
}

export interface ContractTemplateRetrieveByIdRequest extends ContractTemplateBaseRequest {}

export interface ContractTemplateApproveRequest extends ContractTemplateBaseRequest {
  decision_notes?: string[]
}

export interface ContractTemplateRejectRequest extends ContractTemplateBaseRequest {
  /** Reason for rejecting the contract template */
  reason: string
}
