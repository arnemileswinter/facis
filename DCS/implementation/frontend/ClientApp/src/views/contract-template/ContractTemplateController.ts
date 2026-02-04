import { ref, onMounted } from 'vue'
import { ContractTemplateService } from '../../services/contract-template-service'

export function useTemplateTable() {
    const templates = ref([])
    const loading = ref(true)
    const error = ref('')

    const loadTemplates = async () => {
        loading.value = true
        error.value = ''
        try {
            const data = await ContractTemplateService.retrieve()
            console.log(data)

        } catch (err: any) {
            error.value = err.message || 'Fehler beim Laden der Templates'
        } finally {
            loading.value = false
        }
    }

    const refresh = () => loadTemplates()  // Für manuelles Refresh

    const getTemplateById = async (id: string) => {
        try {
            return await ContractTemplateService.retrieveById(id)
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
