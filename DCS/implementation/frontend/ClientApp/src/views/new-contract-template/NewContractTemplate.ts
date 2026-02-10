import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { ContractTemplate } from '../../models/contract-template'
import { ContractTemplateService } from '../../services/contract-template-service'

export function useContractTemplateController() {
    const route = useRoute()
    const router = useRouter()

    const isSubmitting = ref(false)
    const isEditMode = computed(() => !!route.params.did)

    const form = ref<Partial<ContractTemplate>>({
        name: '',
        description: '',
        state: 'DRAFT' as any,
        version: 1,
        meta_data: {}
    })

    onMounted(async () => {
        if (isEditMode.value) {
            await fetchTemplate(route.params.did as string)
        }
    })

    const fetchTemplate = async (did: string) => {
            return await ContractTemplateService.retrieveById(did)
    }

    const submit = async () => {
        isSubmitting.value = true
        try {
            if (isEditMode.value) {
                console.log('Update an API:', form.value)
            } else {
                return await ContractTemplateService.create(form)
            }
            router.push('/')
        } finally {
            isSubmitting.value = false
        }
    }

    return {
        form,
        isEditMode,
        isSubmitting,
        submit,
        cancel: () => router.back()
    }
}