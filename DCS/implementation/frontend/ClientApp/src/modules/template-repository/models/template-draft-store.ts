import type { DocumentOutline, DocumentBlock, SemanticCondition, MetaData, DocumentTypeValue } from "@template-repository/models/contract-templace"

interface TemplateDraftState {
  did: string | null
  documentOutline: DocumentOutline
  documentBlocks: DocumentBlock[]
  semanticConditions: SemanticCondition[]
  customMetaData: MetaData[]
  type: DocumentTypeValue
}



export type { TemplateDraftState }