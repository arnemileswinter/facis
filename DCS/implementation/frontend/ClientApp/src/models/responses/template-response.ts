import type { ContractTemplateType } from "../../types/contract-template-type"
import type { ContractTemplate } from "../contract-template"

interface ContractTemplateBaseResponse {
  did: string
}

interface ContractTemplateBaseRetrieveResponse extends ContractTemplateBaseResponse {
  document_number: number
  version: number
  state: ContractTemplateType
  name?: string
  description?: string
  created_at: string
  updated_at: string
  /** The template data of the contract template */
  template_data: any
}

export interface ContractTemplateResponse extends ContractTemplateBaseResponse {}

export interface ContractTemplateSubmitResponse extends ContractTemplateBaseResponse {}

export interface ContractTemplateUpdateResponse extends ContractTemplateBaseResponse {}

export interface ContractTemplateRetrieveResponse {
  contract_templates: ContractTemplate[]
  review_tasks: any[]
  approval_tasks: any[]
}

export interface ContractTemplateRetrieveByIdResponse extends ContractTemplateBaseRetrieveResponse {
  created_by: string
}

export interface ContractTempleSearchResponse extends ContractTemplateBaseResponse {
  document_number: number
  version: number
  state: ContractTemplateType
  name?: string
  description?: string
  created_at: string
  updated_at: string
}

export interface ContractTemplateApproveResponse extends ContractTemplateBaseResponse {}

export interface ContractTemplateRejectResponse extends ContractTemplateBaseResponse {}
