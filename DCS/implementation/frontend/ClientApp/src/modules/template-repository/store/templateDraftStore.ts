import { defineStore } from 'pinia'
import type { TemplateDraftState } from "@template-repository/models/template-draft-store";

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
    reset() {
      this.$reset()
    }
  }
})

/** Returns a copy of defaultState so store state does not share 
 *  refs with defaultState; mutations in the store do not pollute 
 *  defaultState and $reset() restores correctly.
 **/
function getInitialState(): TemplateDraftState {
  return {
    ...defaultState,
    documentOutline: [...defaultState.documentOutline],
    documentBlocks: [...defaultState.documentBlocks],
    semanticConditions: [...defaultState.semanticConditions],
    customMetaData: [...defaultState.customMetaData],
  }
}