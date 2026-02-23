import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ContractTemplateService } from '../../services/contract-template-service'
import type { ContractTemplateRetrieveByIdRequest } from '../../models/requests/template-request';

interface Clause {
    title: string;
    description: string;
    rules?: SemanticRule[]
}

interface SemanticRule {
    label: string;
    type: string;
    required: boolean;
}

export function useContractTemplateController(did?: string, document_number?: number, version?: number) {
    const router = useRouter()

    const isSubmitting = ref(false)
    const isLoadingSuggestions = ref(false)
    const isEditMode = computed(() => !!did)

    const form = ref({
        name: '',
        description: '',
        clauses: [] as Clause[],
        semantic_rules: [] as SemanticRule[],
        state: 'DRAFT',
        version: 1
    })

    const newClause = ref<Clause>({ title: '', description: '', rules: [] })
    const newRule = ref<SemanticRule>({ label: '', type: '', required: false, })

    const suggestions = ref<SemanticRule[]>([])

    const fetchSuggestions = async () => {
        isLoadingSuggestions.value = true
        try {
            await new Promise(resolve => setTimeout(resolve, 1000))

            suggestions.value = [
                { label: 'Availability_SLA', type: 'text', required: true },
                { label: 'startDate', type: 'date', required: true },
                { label: 'endDate', type: 'date', required: true },
                { label: 'payment_fee', type: 'decimal', required: false },
            ]
        } finally {
            isLoadingSuggestions.value = false
        }
    }

    const addClause = () => {
        if (!newClause.value.title || !newClause.value.description) return

        form.value.clauses.push({
            ...newClause.value,
            rules: [...(newClause.value.rules || [])],
        })

        newClause.value = { title: '', description: '', rules: [] }
    }

    const removeClause = (index: number) => {
        form.value.clauses.splice(index, 1)
    }

    const addRuleToNewClause = (rule: { label: string; type: string; required: boolean }) => {
        newClause.value.rules ||= []
        const exists = newClause.value.rules.some(r => r.label === rule.label)
        if (!exists) newClause.value.rules.push({ ...rule })
    }

    const removeRuleFromNewClause = (idx: number) => {
        newClause.value.rules?.splice(idx, 1)
    }

    const addRuleFromSuggestion = (rule: SemanticRule) => {
        const exists = form.value.semantic_rules.some(r => r.label === rule.label)
        if (!exists) {
            form.value.semantic_rules.push({ ...rule })
        }
    }

    const addNewCustomRule = () => {
        if (!newRule.value.label) return

        form.value.semantic_rules = form.value.semantic_rules || []
        form.value.semantic_rules.push({
            label: newRule.value.label,
            type: newRule.value.type,
            required: newRule.value.required ?? true,
        })

        newRule.value.label = ''
        newRule.value.type = 'Text'
        newRule.value.required = false
    }

    const removeRule = (index: number) => {
        form.value.semantic_rules.splice(index, 1)
    }

    const submit = async () => {
        isSubmitting.value = true
        try {
            console.log("Publishing Template to Repository...", JSON.parse(JSON.stringify(form.value)))
            await new Promise(resolve => setTimeout(resolve, 1500))
            router.push({ name: 'templates.list' })
        } catch (error) {
            console.error("Submission failed", error)
        } finally {
            isSubmitting.value = false
        }
    }

    const retrieveById = async () => {
        if (!did || !document_number ||!version) return

        const request: ContractTemplateRetrieveByIdRequest = {
            did,
            document_number,
            version
        }
        const response = await ContractTemplateService.retrieveById(request)
        if (response) {
            form.value.name = response.name ?? ''
            form.value.description = response.description ?? ''
            form.value.state = response.state
            form.value.version = response.version
        }
    }

    onMounted(async () => {
        fetchSuggestions()
        if (isEditMode.value) {
            await retrieveById()
        }
    })

    return {
        form,
        isEditMode,
        isSubmitting,
        isLoadingSuggestions,
        newClause,
        newRule,
        suggestions,
        addRuleToNewClause,
        removeRuleFromNewClause, 
        addClause,
        removeClause,
        addRuleFromSuggestion,
        addNewCustomRule,
        removeRule,
        submit,
        cancel: () => router.back()
    }
}
