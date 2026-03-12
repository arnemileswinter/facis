import type { PartialContractTemplate } from '@/models/contract-template'
import { defineStore } from 'pinia'
import { ref, type Ref } from 'vue'

export const useContractTemplatesStore = defineStore('contractTemplates', () => {
  const contractTemplates: Ref<PartialContractTemplate[]> = ref([])

  return { contractTemplates }
})
