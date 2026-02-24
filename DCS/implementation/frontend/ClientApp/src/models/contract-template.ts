import type { ContractTemplateType } from "@/types/contract-template-type"
import type { TemplateType } from "@/types/template-type"

export interface ContractTemplate {
    did: string
    created_by: string
    created_at: string
    document_number: number
    version: number
    template_type?: TemplateType
    state: ContractTemplateType
    name?: string
    description?: string
    template_data?: any
    updated_at: string
    clauses?: []
}
