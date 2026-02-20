import type {
  DocumentOutline,
  DocumentBlock,
  SemanticCondition,
  MetaData,
  DocumentTypeValue,
  DocumentBlockType,
} from "@template-repository/models/contract-templace"

interface TemplateDraftState {
  did: string | null
  documentOutline: DocumentOutline
  documentBlocks: DocumentBlock[]
  semanticConditions: SemanticCondition[]
  customMetaData: MetaData[]
  type: DocumentTypeValue
}

/** Payload for adding a new block. */
export interface AddBlockPayload {
  blockType: DocumentBlockType
  text: string
  title?: string
  conditionId?: string
  templateId?: string
}

export type { TemplateDraftState }