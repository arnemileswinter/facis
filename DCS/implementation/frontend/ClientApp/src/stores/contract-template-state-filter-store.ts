import type { ContractTemplateType } from '@/types/contract-template-type'
import { defineStore } from 'pinia'
import { ref, type Ref } from 'vue'

export const useContractTemplateStateFilterStore = defineStore('contractTemplateStateFilter', () => {
  const stateFilters: Ref<Set<ContractTemplateType>> = ref(new Set())

  function setFilter(type: ContractTemplateType) {
    stateFilters.value.add(type)
  }

  function removeFilter(type: ContractTemplateType) {
    stateFilters.value.delete(type)
  }

  return { stateFilters, setFilter, removeFilter }
})
