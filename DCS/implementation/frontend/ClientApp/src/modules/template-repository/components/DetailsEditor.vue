<template>
    <div class="grid grid-cols-1 gap-4">
        <!-- Contract Kind -->
        <fieldset class="fieldset p-0 border-none">
            <legend class="fieldset-legend">Contract Type</legend>
            <div class="grid grid-cols-2 gap-3 mt-1">
                <div class="card border-2 transition-all pointer-events-none"
                    :class="templateType === TemplateType.frameContract
                        ? 'border-primary bg-primary/5'
                        : 'border-base-300'">
                    <div class="card-body p-4 gap-1">
                        <span class="card-title text-sm">Frame Contract</span>
                        <p class="text-xs text-base-content/60 font-normal">Top-level agreement that groups subcontracts</p>
                    </div>
                </div>
                <div class="card border-2 transition-all pointer-events-none"
                    :class="templateType === TemplateType.subContract
                        ? 'border-primary bg-primary/5'
                        : 'border-base-300'">
                    <div class="card-body p-4 gap-1">
                        <span class="card-title text-sm">Subcontract</span>
                        <p class="text-xs text-base-content/60 font-normal">Scoped agreement under a frame contract</p>
                    </div>
                </div>
            </div>
        </fieldset>

        <fieldset class="fieldset p-0 border-none">
            <legend class="fieldset-legend">Global Name</legend>
            <input v-model="name" class="input input-bordered w-full" type="text" required />
        </fieldset>

        <fieldset class="fieldset p-0 border-none">
            <legend class="fieldset-legend">Base Description</legend>
            <textarea v-model="description" class="textarea textarea-bordered w-full h-24" required></textarea>
        </fieldset>

        <!-- Subcontracts (only for frame contracts) -->
        <fieldset v-if="templateType === TemplateType.frameContract" class="fieldset p-0 border-none">
            <legend class="fieldset-legend cursor-pointer select-none inline-flex items-center gap-1.5"
                @click="showSubcontractPicker = !showSubcontractPicker">
                Subcontract Templates
                <svg xmlns="http://www.w3.org/2000/svg"
                    class="w-3 h-3 transition-transform duration-200 opacity-60"
                    :class="{ 'rotate-180': showSubcontractPicker }"
                    fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
            </legend>

            <!-- Collapsible picker -->
            <div v-show="showSubcontractPicker" class="mt-1">
                <input v-model="subcontractSearchQuery"
                    class="input input-bordered input-sm w-full"
                    placeholder="Search templates…" />

                <ul class="menu menu-sm w-full bg-base-200 rounded-box mt-1 max-h-48 overflow-y-auto flex-nowrap">
                    <li v-if="!filteredSubcontractTemplates.length">
                        <span class="text-base-content/40 italic text-xs pointer-events-none">
                            {{ subcontractSearchQuery ? 'No results' : 'All templates already selected' }}
                        </span>
                    </li>
                    <li v-for="t in filteredSubcontractTemplates" :key="t.did">
                        <button type="button" @click="addSubcontractTemplate(t)"
                            class="group flex flex-col items-start gap-0">
                            <span class="font-medium text-sm">{{ t.name }}</span>
                            <span class="text-xs text-base-content/50 italic overflow-hidden max-h-0 group-hover:max-h-12 transition-all duration-200 ease-in-out">
                                {{ t.description }}
                            </span>
                        </button>
                    </li>
                </ul>
            </div>

            <!-- Selected templates (always visible) -->
            <div v-if="selectedSubcontractDids.length" class="flex flex-wrap gap-2 mt-3">
                <div v-for="did in selectedSubcontractDids" :key="did"
                    class="badge badge-primary badge-outline gap-1 py-3">
                    <span>{{ getSubcontractTemplateName(did) }}</span>
                    <button type="button" @click="removeSubcontractTemplate(did)"
                        class="text-error hover:opacity-70 transition-opacity">✕</button>
                </div>
            </div>
            <p v-else class="fieldset-label mt-2">No subcontract templates selected yet.</p>
        </fieldset>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute } from 'vue-router'
import { useTemplateDraftStore } from '@template-repository/store/templateDraftStore'
import { TemplateType } from '@template-repository/models/contract-templace'
import { ContractTemplateService } from '@/services/contract-template-service'

interface SubcontractTemplate {
    did: string
    name: string
    description: string
}

const store = useTemplateDraftStore()
const { templateType } = storeToRefs(store)

const name = computed({
  get: () => store.name,
  set: (value: string) => store.updateName(value.trim())
})

const description = computed({
  get: () => store.description,
  set: (value: string) => store.updateDescription(value)
})

const route = useRoute()

const selectedSubcontractDids = ref<string[]>([])
const showSubcontractPicker = ref(false)
const subcontractSearchQuery = ref('')

const availableSubcontractTemplates = ref<SubcontractTemplate[]>([
    { did: 'did:facis:tmpl:sc:001', name: 'IT Services Agreement', description: 'General IT services including consulting, implementation, and technical support.' },
    { did: 'did:facis:tmpl:sc:002', name: 'Hardware Procurement', description: 'Purchase and delivery of physical hardware components and equipment.' },
    { did: 'did:facis:tmpl:sc:003', name: 'Software License', description: 'Licensing terms for proprietary or third-party software products.' },
    { did: 'did:facis:tmpl:sc:004', name: 'Maintenance & Support', description: 'Ongoing maintenance, updates, and technical support services.' },
    { did: 'did:facis:tmpl:sc:005', name: 'Cloud Infrastructure Services', description: 'Provisioning and management of cloud-based infrastructure resources.' },
    { did: 'did:facis:tmpl:sc:006', name: 'Data Processing Agreement', description: 'GDPR-compliant data processing terms between controller and processor.' },
    { did: 'did:facis:tmpl:sc:007', name: 'Consulting Services', description: 'Professional advisory and strategic consulting engagements.' },
])

const filteredSubcontractTemplates = computed(() => {
    const q = subcontractSearchQuery.value.toLowerCase()
    return availableSubcontractTemplates.value.filter(t =>
        !selectedSubcontractDids.value.includes(t.did) &&
        (q === '' || t.name.toLowerCase().includes(q) || t.did.toLowerCase().includes(q))
    )
})

const getSubcontractTemplateName = (did: string) =>
    availableSubcontractTemplates.value.find(t => t.did === did)?.name ?? did

const addSubcontractTemplate = (template: SubcontractTemplate) => {
    if (!selectedSubcontractDids.value.includes(template.did)) {
        selectedSubcontractDids.value.push(template.did)
    }
    subcontractSearchQuery.value = ''
}

const removeSubcontractTemplate = (did: string) => {
    const idx = selectedSubcontractDids.value.indexOf(did)
    if (idx !== -1) selectedSubcontractDids.value.splice(idx, 1)
}

onMounted(async () => {
    const did = route.params.did
    const documentNumber = route.query.document_number
    const version = route.query.version

    if (
        did && documentNumber && version &&
        !Array.isArray(did) && !Array.isArray(documentNumber) && !Array.isArray(version)
    ) {
        const response = await ContractTemplateService.retrieveById({
            did,
            document_number: parseInt(documentNumber),
            version: parseInt(version),
        })
        if (response) {
            name.value = response.name ?? ''
            description.value = response.description ?? ''
        }
    }
})
</script>
