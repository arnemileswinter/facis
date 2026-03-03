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
  name: string
  description: string
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
  // #### For Clause ####
  clauseBlockId?: string
  conditionIds?: string[]
  // #### For ApprovedTemplate ####
  templateId?: string
}

export interface AddBlockOptions {
  addToOutline?: boolean
}

export type { TemplateDraftState }