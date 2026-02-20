import { defineStore } from 'pinia'
import type { TemplateDraftState, AddBlockPayload } from "@template-repository/models/template-draft-store"
import type { DocumentOutlineBlock, DocumentBlock } from "@template-repository/models/contract-templace"
import { DocumentBlockType, DocumentType } from "@template-repository/models/contract-templace"

const storeId = "templateDraft"
const defaultState: Readonly<TemplateDraftState> = {
  did: null,
  documentOutline: [],
  documentBlocks: [],
  semanticConditions: [],
  customMetaData: [],
  type: 'subContract',
}

export const useTemplateDraftStore = defineStore(storeId, {
  state: (): TemplateDraftState => getInitialState(),
  getters: {
    hasTemplateId(): boolean { return !!this.did }
  },
  actions: {
    /**
     * Adds a new block under the given parent at the given index.
     * 
     * subContract: cannot add APPROVED_TEMPLATE. frameContract: can only add APPROVED_TEMPLATE.
     * 
     * @param parentBlockId - blockId of the outline node (parent) under which to insert
     * @param insertIndex - index in the parent's children array (0 = first)
     * @returns The new block's blockId.
     */
    addBlock(parentBlockId: string, insertIndex: number, payload: AddBlockPayload): string {
      if (this.type === DocumentType.subContract && payload.blockType === DocumentBlockType.ApprovedTemplate) {
        throw new Error('subContract template cannot add APPROVED_TEMPLATE blocks')
      }
      if (this.type === DocumentType.frameContract && payload.blockType !== DocumentBlockType.ApprovedTemplate) {
        throw new Error('frameContract template can only add APPROVED_TEMPLATE blocks')
      }
      return addBlock(this.documentOutline, this.documentBlocks, parentBlockId, insertIndex, payload)
    },
    /** Removes the block and all its descendants from documentOutline and documentBlocks. */
    deleteBlock(blockId: string): void {
      deleteBlock(this.documentOutline, this.documentBlocks, blockId)
    },
    /** Updates block fields. */
    updateBlock(blockId: string, payload: { title?: string; text?: string }): void {
      const block = this.documentBlocks.find((b) => b.blockId === blockId)
      if (!block) return
      if (payload.title !== undefined) block.title = payload.title
      if (payload.text !== undefined) block.text = payload.text
    },

    reset() {
      Object.assign(this, getInitialState())
    }
  }
})

/** Creates a document outline item */
function createOutlineItem(overrides?: Partial<Pick<DocumentOutlineBlock, 'blockId' | 'isRoot' | 'children'>>): DocumentOutlineBlock {
  return {
    blockId: overrides?.blockId ?? crypto.randomUUID(),
    isRoot: overrides?.isRoot ?? false,
    children: overrides?.children ?? [],
  }
}

function addBlock(
  outline: DocumentOutlineBlock[],
  blocks: DocumentBlock[],
  parentBlockId: string,
  insertIndex: number,
  payload: AddBlockPayload
): string {
  const blockId = crypto.randomUUID()
  const block = createBlockFromPayload(blockId, payload)
  const parent = outline.find((b) => b.blockId === parentBlockId)
  if (!parent) {
    throw new Error(`addBlock: parent not found: ${parentBlockId}`)
  }
  parent.children.splice(insertIndex, 0, blockId)
  if (payload.blockType === DocumentBlockType.Section) {
    outline.push(createOutlineItem({ blockId, isRoot: false, children: [] }))
  }
  blocks.push(block)
  return blockId
}

/**
 * Removes the block and all its descendants from outline and blocks. Mutates both arrays.
 */
function deleteBlock(
  outline: DocumentOutlineBlock[],
  blocks: DocumentBlock[],
  blockId: string
): void {
  const outlineByBlockId = new Map(outline.map((b) => [b.blockId, b]))
  const parent = outline.find((b) => b.children.includes(blockId))
  if (!parent) {
    return
  }
  const toRemove = collectDescendantBlockIds(blockId, outlineByBlockId)
  parent.children = parent.children.filter((id) => id !== blockId)
  const outlineToKeep = outline.filter((b) => b.isRoot === true || !toRemove.has(b.blockId))
  outline.length = 0
  outline.push(...outlineToKeep)
  const blocksToKeep = blocks.filter((b) => !toRemove.has(b.blockId))
  blocks.length = 0
  blocks.push(...blocksToKeep)
}

function createBlockFromPayload(blockId: string, payload: AddBlockPayload): DocumentBlock {
  const text = payload.text ?? ''
  switch (payload.blockType) {
    case DocumentBlockType.Section:
      return { blockId, type: DocumentBlockType.Section, text }
    case DocumentBlockType.Text:
      return { blockId, type: DocumentBlockType.Text, text }
    case DocumentBlockType.Clause:
      return { blockId, type: DocumentBlockType.Clause, text, conditionId: payload.conditionId ?? '' }
    case DocumentBlockType.ApprovedTemplate:
      return { blockId, type: DocumentBlockType.ApprovedTemplate, text, templateId: payload.templateId ?? '' }
    default:
      throw new Error(`Unknown blockType: ${payload.blockType}`)
  }
}

/** Returns a Set of blockId and all descendant block ids in the outline. */
function collectDescendantBlockIds(
  blockId: string,
  outlineByBlockId: Map<string, DocumentOutlineBlock>
): Set<string> {
  const set = new Set<string>([blockId])
  const node = outlineByBlockId.get(blockId)
  const childIds = node?.children ?? []
  childIds.forEach((id) => collectDescendantBlockIds(id, outlineByBlockId).forEach((x) => set.add(x)))
  return set
}

/** Returns a copy of defaultState so store state does not share 
 *  refs with defaultState; mutations in the store do not pollute 
 *  defaultState and $reset() restores correctly.
 **/
function getInitialState(): TemplateDraftState {
  return {
    ...defaultState,
    // Root is not created by the user
    documentOutline: [createOutlineItem({ isRoot: true, children: [] })],
    documentBlocks: [...defaultState.documentBlocks],
    semanticConditions: [...defaultState.semanticConditions],
    customMetaData: [...defaultState.customMetaData],
  }
}