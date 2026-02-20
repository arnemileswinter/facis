import { defineStore } from 'pinia'
import type { TemplateEditorUiState, TemplateEditorTabId } from "@template-repository/models/template-editor-ui-store";

const storeId = "templateDraft"
const defaultState: TemplateEditorUiState = {
  activeTab: 'details',
  tabs: [
    { id: 'details', label: 'Details' },
    { id: 'semantic', label: 'Semantic Rules' },
    { id: 'clauses', label: 'Clauses' },
    { id: 'builder', label: 'Builder' },
    { id: 'meta', label: 'Meta Data' },
  ]
}

export const useTemplateEditorUiStore = defineStore(storeId, {
  state: (): TemplateEditorUiState => defaultState,
  getters: {},
  actions: {
    setActiveTab(tab: TemplateEditorTabId) {
      this.activeTab = tab
    }
  }
})