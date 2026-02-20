export type TemplateEditorTabId = 'details' | 'semantic' | 'clauses' | 'builder' | 'meta'

/** UI state for template create/edit page */
interface TemplateEditorUiState {
  activeTab: TemplateEditorTabId
  tabs: [
    { id: 'details', label: string },
    { id: 'semantic', label: string },
    { id: 'clauses', label: string },
    { id: 'builder', label: string },
    { id: 'meta', label: string },
  ]
}

export type { TemplateEditorUiState }