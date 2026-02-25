import type { ContractTemplateState } from "@/types/contract-template-state"
import type { TemplateType } from "@/types/template-type"

export interface ContractTemplate {
    did: string
    created_by: string
    created_at: string
    document_number: number
    version: number
    template_type?: TemplateType
    state: ContractTemplateState
    name?: string
    description?: string
    template_data?: any
    updated_at: string
    clauses?: []
}
