import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

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

export function useContractTemplateController() {
    const route = useRoute()
    const router = useRouter()

    const isSubmitting = ref(false)
    const isLoadingSuggestions = ref(false)
    const isEditMode = computed(() => !!route.params.did)

    const form = ref({
        name: '',
        description: '',
        clauses: [] as Clause[],
        semantic_rules: [] as SemanticRule[],
        state: 'DRAFT',
        version: 1
    })

    const newClause = ref<Clause>({ title: '', description: '' })
    const newRule = ref<SemanticRule>({ label: '', type: '', required: false })

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
        form.value.clauses.push({ ...newClause.value })
        newClause.value = { title: '', description: '' }
    }

    const removeClause = (index: number) => {
        form.value.clauses.splice(index, 1)
    }

    const addRuleFromSuggestion = (rule: SemanticRule) => {
        const exists = form.value.semantic_rules.some(r => r.label === rule.label)
        if (!exists) {
            form.value.semantic_rules.push({ ...rule })
        }
    }

    const addNewCustomRule = () => {
        if (!newRule.value.label || !newRule.value.type || !newRule.value.required) return
        form.value.semantic_rules.push({ ...newRule.value })
        newRule.value = { label: '', type: '', required: false }
    }

    const removeRule = (index: number) => {
        form.value.semantic_rules.splice(index, 1)
    }

    const submit = async () => {
        isSubmitting.value = true
        try {
            console.log("Publishing Template to Repository...", JSON.parse(JSON.stringify(form.value)))
            await new Promise(resolve => setTimeout(resolve, 1500))
            router.push('/')
        } catch (error) {
            console.error("Submission failed", error)
        } finally {
            isSubmitting.value = false
        }
    }

    onMounted(() => {
        fetchSuggestions()
        if (isEditMode.value) {
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

        addClause,
        removeClause,
        addRuleFromSuggestion,
        addNewCustomRule,
        removeRule,
        submit,
        cancel: () => router.back()
    }
}