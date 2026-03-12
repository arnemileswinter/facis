import { onMounted, ref, type Ref } from 'vue'
import type { PartialContractTemplate } from '../../models/contract-template'
import { ContractTemplateService } from '../../services/contract-template-service'
import type { ContractTemplateApprovalTask } from '@/models/contract-template-approval-task'
import type { ContractTemplateReviewTask } from '@/models/contract-template-review-task'
import { useAuthStore } from '@/stores/auth-store'

export function useTemplateTable() {
    const templates: Ref<PartialContractTemplate[]> = ref([])
    const approvalTasks: Ref<ContractTemplateApprovalTask[]> = ref([])
    const reviewTasks: Ref<ContractTemplateReviewTask[]> = ref([])
    const loading = ref(true)
    const error = ref('')
    const authStore = useAuthStore()

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

    const getTemplateById = async (did: string) => {
        try {
            return await ContractTemplateService.retrieveById({ did })
        } catch (err: any) {
            console.error('Template konnte nicht geladen werden:', err)
            return null
        }
    }
    onMounted(loadTemplates)

    const hasReviewTask = (template: PartialContractTemplate): boolean => {
        // Bug: authStore.user is not userId, it's a JWT token. We need to decode it to get the userId.
        const currentUser = authStore.user
        if (!currentUser) return false
        return reviewTasks.value.some((task) => {
            const isDidMatch = task.did === template.did
            const isVersionMatch = !template.version || task.version === template.version
            const isDocumentNumberMatch = !template.document_number || task.document_number === template.document_number
            return (
                isDidMatch &&
                isVersionMatch &&
                isDocumentNumberMatch &&
                task.reviewer === currentUser
            )
        })
    }

    const hasApprovalTask = (template: PartialContractTemplate): boolean => {
        const currentUser = authStore.user
        if (!currentUser) return false
        return approvalTasks.value.some((task) => {
            const isDidMatch = task.did === template.did
            const isVersionMatch = !template.version || task.version === template.version
            const isDocumentNumberMatch = !template.document_number || task.document_number === template.document_number
            return (
                isDidMatch &&
                isVersionMatch &&
                isDocumentNumberMatch &&
                task.approver === currentUser
            )
        })
    }

    return {
        templates,
        approvalTasks,
        reviewTasks,
        loading,
        error,
        loadTemplates,
        refresh,
        getTemplateById,
        hasReviewTask,
        hasApprovalTask,
    }
}
