import { onMounted, ref, type Ref } from 'vue'
import type { PartialContractTemplate } from '../../models/contract-template'
import { ContractTemplateService } from '../../services/contract-template-service'
import type { ContractTemplateApprovalTask } from '@/models/contract-template-approval-task'
import type { ContractTemplateReviewTask } from '@/models/contract-template-review-task'

export function useTemplateTable() {
    const templates: Ref<PartialContractTemplate[]> = ref([])
    const approvalTasks: Ref<ContractTemplateApprovalTask[]> = ref([])
    const reviewTasks: Ref<ContractTemplateReviewTask[]> = ref([])
    const loading = ref(true)
    const error = ref('')

    const loadTemplates = async () => {
        loading.value = true
        error.value = ''
        try {
            const data = await ContractTemplateService.retrieve()
            console.log(data)
            templates.value = data.contract_templates
            approvalTasks.value = data.approval_tasks
            reviewTasks.value = data.review_tasks

        } catch (err: any) {
            error.value = err.message || 'Fehler beim Laden der Templates'
        } finally {
            loading.value = false
        }
    }

    const refresh = () => loadTemplates()  // Für manuelles Refresh

    const getTemplateById = async (did: string, version: number, document_number: string) => {
        try {
            return await ContractTemplateService.retrieveById({ did, version, document_number })
        } catch (err: any) {
            console.error('Template konnte nicht geladen werden:', err)
            return null
        }
    }
    onMounted(loadTemplates)

    return {
        templates,
        approvalTasks,
        reviewTasks,
        loading,
        error,
        loadTemplates,
        refresh,
        getTemplateById
    }
}
