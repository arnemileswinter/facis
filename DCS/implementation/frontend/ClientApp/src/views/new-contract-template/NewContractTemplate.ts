import type { ContractTemplate } from '@/models/contract-template';
import type { ContractTemplateRetrieveByIdRequest } from '@/models/requests/template-request';
import { ContractTemplateService } from '@/services/contract-template-service';
import { TemplateType, type TemplateTypeValue } from "@template-repository/models/contract-templace";
import { computed, onMounted, ref, type Ref } from 'vue';
import { useRouter } from 'vue-router';
interface SubcontractTemplate {
    did: string;
    name: string;
    description: string;
}

export function useContractTemplateController(did?: string, document_number?: number, version?: number) {
    const router = useRouter()

    const contractTemplate: Ref<ContractTemplate | null> = ref(null)

    const isSubmitting = ref(false)
    const isEditMode = computed(() => !!did)
    const form = computed(() => {
        return {
            name: contractTemplate.value?.name ?? '',
            description: contractTemplate.value?.description ?? '',
            contract_kind: TemplateType.subContract as TemplateTypeValue,
            subcontract_template_dids: [] as string[],
            state: contractTemplate.value?.state ?? 'DRAFT',
            version: contractTemplate.value?.state ?? 1
        }
    })

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
            contractTemplate.value = response
        }
    }

    onMounted(async () => {
        if (isEditMode.value) {
            await retrieveById()
        }
    })

    return {
        form,
        isEditMode,
        isSubmitting,
        subcontractSearchQuery,
        filteredSubcontractTemplates,
        getSubcontractTemplateName,
        addSubcontractTemplate,
        removeSubcontractTemplate,
        submit,
        cancel: () => router.back(),
    }
}
