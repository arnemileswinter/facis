import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

interface Clause {
    title: string;
    content: string;
}

interface SemanticRule {
    attr: string;
    op: string;
    val: string;
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

    const newClause = ref<Clause>({ title: '', content: '' })
    const newRule = ref<SemanticRule>({ attr: '', op: 'MIN', val: '' })
    
    const suggestions = ref<SemanticRule[]>([])

    const fetchSuggestions = async () => {
        isLoadingSuggestions.value = true
        try {
            await new Promise(resolve => setTimeout(resolve, 1000))
            
            suggestions.value = [
                { attr: 'Availability_SLA', op: 'MIN', val: '99.9%' },
                { attr: 'Liability_Limit', op: 'MAX', val: '1.000.000€' },
                { attr: 'Data_Location', op: 'EQUALS', val: 'EU-Only' },
                { attr: 'Support_Response', op: 'MAX', val: '2h' }
            ]
        } finally {
            isLoadingSuggestions.value = false
        }
    }

    const addClause = () => {
        if (!newClause.value.title || !newClause.value.content) return
        form.value.clauses.push({ ...newClause.value })
        newClause.value = { title: '', content: '' }
    }

    const removeClause = (index: number) => {
        form.value.clauses.splice(index, 1)
    }

    const addRuleFromSuggestion = (rule: SemanticRule) => {
        const exists = form.value.semantic_rules.some(r => r.attr === rule.attr)
        if (!exists) {
            form.value.semantic_rules.push({ ...rule })
        }
    }

    const addNewCustomRule = () => {
        if (!newRule.value.attr || !newRule.value.val) return
        form.value.semantic_rules.push({ ...newRule.value })
        newRule.value = { attr: '', op: 'MIN', val: '' }
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