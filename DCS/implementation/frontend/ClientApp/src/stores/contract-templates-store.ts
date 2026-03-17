import type { PartialContractTemplate } from '@/models/contract-template'
import type { ContractTemplateApprovalTask } from '@/models/contract-template-approval-task'
import type { ContractTemplateReviewTask } from '@/models/contract-template-review-task'
import { defineStore } from 'pinia'
import { ref, type Ref } from 'vue'

export const useContractTemplatesStore = defineStore('contractTemplates', () => {
  const contractTemplates: Ref<PartialContractTemplate[]> = ref([])
  const reviewTasks: Ref<ContractTemplateReviewTask[]> = ref([])
  const approvalTasks: Ref<ContractTemplateApprovalTask[]> = ref([])

  return { contractTemplates, reviewTasks, approvalTasks }
})
