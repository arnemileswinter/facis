import type { ContractTemplateType } from "../types/contract-template-type"

export interface ContractTemplate {
    did: string
    created_by: string
    created_at: string
    document_number: number
    version: number
    state: ContractTemplateType
    name: string
    description?: string
    meta_data?: any
    clauses?: []
}