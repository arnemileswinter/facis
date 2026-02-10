interface ContractTemplateBaseResponse {
  did: string
}

interface ContractTemplateBaseRetrieveResponse extends ContractTemplateBaseResponse {
  document_number: number
  version: number
  state: string
  name?: string
  description?: string
  created_by: string
  created_at: string
  /** The metadata of the contract template */
  meta_data: any
}

export interface ContractTemplateResponse extends ContractTemplateBaseResponse {}

export interface ContractTemplateSubmitResponse extends ContractTemplateBaseResponse {}

export interface ContractTemplateUpdateResponse extends ContractTemplateBaseResponse {}

export interface ContractTemplateRetrieveResponse extends ContractTemplateBaseRetrieveResponse {}

export interface ContractTemplateRetrieveByIdResponse extends ContractTemplateBaseRetrieveResponse {}

export interface ContractTemplateApproveResponse extends ContractTemplateBaseResponse {}

export interface ContractTemplateRejectResponse extends ContractTemplateBaseResponse {}
