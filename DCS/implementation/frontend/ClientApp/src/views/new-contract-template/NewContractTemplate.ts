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

interface SubcontractTemplate {
    did: string;
    name: string;
    description: string;
}

export function useContractTemplateController(did?: string, document_number?: number, version?: number) {
    const router = useRouter()

    const isSubmitting = ref(false)
    const isLoadingSuggestions = ref(false)
    const isEditMode = computed(() => !!did)
    const selectedRules = ref([])
    const form = ref({
        name: '',
        description: '',
        contract_kind: 'frame_contract' as 'frame_contract' | 'subcontract',
        subcontract_template_dids: [] as string[],
        clauses: [] as Clause[],
        semantic_rules: [] as SemanticRule[],
        state: 'DRAFT',
        version: 1
    })

    const newClause = ref<Clause>({ title: '', description: '', rules: [] })
    const newRule = ref<SemanticRule>({ label: '', type: '', required: false, })

    // Subcontract template picker
    const availableSubcontractTemplates = ref<SubcontractTemplate[]>([
        { did: 'did:facis:tmpl:sc:001', name: 'IT Services Agreement', description: 'General IT services including consulting, implementation, and technical support.' },
        { did: 'did:facis:tmpl:sc:002', name: 'Hardware Procurement', description: 'Purchase and delivery of physical hardware components and equipment.' },
        { did: 'did:facis:tmpl:sc:003', name: 'Software License', description: 'Licensing terms for proprietary or third-party software products.' },
        { did: 'did:facis:tmpl:sc:004', name: 'Maintenance & Support', description: 'Ongoing maintenance, updates, and technical support services.' },
        { did: 'did:facis:tmpl:sc:005', name: 'Cloud Infrastructure Services', description: 'Provisioning and management of cloud-based infrastructure resources.' },
        { did: 'did:facis:tmpl:sc:006', name: 'Data Processing Agreement', description: 'GDPR-compliant data processing terms between controller and processor.' },
        { did: 'did:facis:tmpl:sc:007', name: 'Consulting Services', description: 'Professional advisory and strategic consulting engagements.' },
    ])

    const subcontractSearchQuery = ref('')

    const filteredSubcontractTemplates = computed(() => {
        const q = subcontractSearchQuery.value.toLowerCase()
        return availableSubcontractTemplates.value.filter(t =>
            !form.value.subcontract_template_dids.includes(t.did) &&
            (q === '' || t.name.toLowerCase().includes(q) || t.did.toLowerCase().includes(q))
        )
    })

    const getSubcontractTemplateName = (did: string) =>
        availableSubcontractTemplates.value.find(t => t.did === did)?.name ?? did

    const addSubcontractTemplate = (template: SubcontractTemplate) => {
        if (!form.value.subcontract_template_dids.includes(template.did)) {
            form.value.subcontract_template_dids.push(template.did)
        }
        subcontractSearchQuery.value = ''
    }

    const removeSubcontractTemplate = (did: string) => {
        const idx = form.value.subcontract_template_dids.indexOf(did)
        if (idx !== -1) form.value.subcontract_template_dids.splice(idx, 1)
    }

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
            rules: [...selectedRules.value],
        })
        newClause.value = {
            title: '',
            description: '',
            rules: []
        }
        selectedRules.value = []
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
        if (!did || !document_number || !version) return

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
        subcontractSearchQuery,
        filteredSubcontractTemplates,
        getSubcontractTemplateName,
        addSubcontractTemplate,
        removeSubcontractTemplate,
        suggestions,
        addRuleToNewClause,
        removeRuleFromNewClause,
        addClause,
        removeClause,
        addRuleFromSuggestion,
        addNewCustomRule,
        removeRule,
        submit,
        selectedRules,
        cancel: () => router.back()
    }
}
