import { onMounted, ref, type Ref } from 'vue'
import type { ContractTemplate } from '../../models/contract-template'
import { ContractTemplateService } from '../../services/contract-template-service'

export function useTemplateTable() {
    const templates: Ref<ContractTemplate[]> = ref([])
    const loading = ref(true)
    const error = ref('')

    const loadTemplates = async () => {
        loading.value = true
        error.value = ''
        try {
            const data = await ContractTemplateService.retrieve()
            console.log(data)
            templates.value = data

        } catch (err: any) {
            error.value = err.message || 'Fehler beim Laden der Templates'
        } finally {
            loading.value = false
        }
    }

    const refresh = () => loadTemplates()  // Für manuelles Refresh

    const getTemplateById = async (id: string) => {
        try {
            return await ContractTemplateService.retrieveById({
                did: id,
                document_number: 1,
                version: 1
            })
        } catch (err: any) {
            console.error('Template konnte nicht geladen werden:', err)
            return null
        }
    }
    onMounted(loadTemplates)

    return {
        templates,
        loading,
        error,
        loadTemplates,
        refresh,
        getTemplateById
    }
}
