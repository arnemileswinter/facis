export type TemplateEditorTabId = 'details' | 'semantic' | 'clauses' | 'builder' | 'meta'
export interface AddBlockModalContext {
  parentBlockId: string
  /** Index in the parent's children array where the new block will be inserted */
  insertIndex: number
}

/** UI state for template create/edit page */
interface TemplateEditorUiState {
  activeTab: TemplateEditorTabId
  tabs: [
    { id: 'details', label: string },
    { id: 'semantic', label: string },
    { id: 'clauses', label: string },
    { id: 'builder', label: string },
    { id: 'meta', label: string },
  ],
  /**
   * When non-null: add-block modal is open
   */
  addBlockModalContext: AddBlockModalContext | null
}

export type { TemplateEditorUiState }