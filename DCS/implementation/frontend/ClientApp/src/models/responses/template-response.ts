import type { ContractTemplateType } from "../../types/contract-template-type"
import type { ContractTemplate } from "../contract-template"

interface ContractTemplateBaseResponse {
  did: string
  document_number: number
  version: number
}

export interface ContractTemplateCreateResponse extends ContractTemplateBaseResponse {}

export interface ContractTemplateSubmitResponse extends ContractTemplateBaseResponse {}

export interface ContractTemplateUpdateResponse extends ContractTemplateBaseResponse {}

export interface ContractTemplateSearchResponse extends ContractTemplateBaseResponse {
  state: ContractTemplateType
  name?: string
  description?: string
  created_at: string
  updated_at: string
}

export interface ContractTemplateRetrieveResponse {
  contract_templates: ContractTemplate[]
  review_tasks: any[]
  approval_tasks: any[]
}

export interface ContractTemplateRetrieveByIdResponse extends ContractTemplateBaseResponse {
  state: ContractTemplateType
  name?: string
  description?: string
  created_by: string
  created_at: string
  updated_at: string
  /** The template data of the contract template */
  template_data: any
}

export interface ContractTemplateApproveResponse extends ContractTemplateBaseResponse {}

export interface ContractTemplateRejectResponse extends ContractTemplateBaseResponse {}
