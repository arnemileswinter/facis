import { defineStore } from 'pinia'
import type {
  TemplateEditorUiState,
  TemplateEditorTabId,
  BlockMovementPreview,
} from '@template-repository/models/template-editor-ui-store'

const storeId = 'templateEditorUi'
const defaultState: Readonly<TemplateEditorUiState> = {
  activeTab: 'details',
  tabs: [
    { id: 'details', label: 'Details' },
    { id: 'semantic', label: 'Semantic Rules' },
    { id: 'clauses', label: 'Clauses' },
    { id: 'builder', label: 'Builder' },
    { id: 'meta', label: 'Meta Data' },
  ],
  addBlockModalContext: null,
  blockMovementPreview: null,
  selectedBlockId: null,
}

export const useTemplateEditorUiStore = defineStore(storeId, {
  state: (): TemplateEditorUiState => getInitialState(),
  getters: {},
  actions: {
    setActiveTab(tab: TemplateEditorTabId) {
      this.activeTab = tab
    },
    openAddBlockModal(parentBlockId: string, insertIndex: number) {
      this.addBlockModalContext = { parentBlockId, insertIndex }
    },
    closeAddBlockModal() {
      this.addBlockModalContext = null
    },
    setBlockMovementPreview(value: BlockMovementPreview | null) {
      this.blockMovementPreview = value
    },
    setSelectedBlockId(blockId: string | null) {
      this.selectedBlockId = blockId
    },
    reset() {
      Object.assign(this, getInitialState())
    }
  }
})

function getInitialState(): TemplateEditorUiState {
  return {
    ...defaultState,
    tabs: [...defaultState.tabs]
  }
}