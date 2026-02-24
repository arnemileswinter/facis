import type {
  DocumentOutline,
  DocumentBlock,
  SemanticCondition,
  MetaData,
  TemplateTypeValue,
  DocumentBlockType,
} from "@template-repository/models/contract-templace"

interface TemplateDraftState {
  did: string | null
  documentOutline: DocumentOutline
  documentBlocks: DocumentBlock[]
  semanticConditions: SemanticCondition[]
  customMetaData: MetaData[]
  templateType: TemplateTypeValue
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