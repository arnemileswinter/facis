import type { ContractTemplateType } from "../types/contract-template-type"

export interface ContractTemplate {
    did: string
    createdBy: string
    createdAt: string
    documentNumber: number
    version: number
    state: ContractTemplateType
    name: string
    description?: string
    meta_data?: any
}