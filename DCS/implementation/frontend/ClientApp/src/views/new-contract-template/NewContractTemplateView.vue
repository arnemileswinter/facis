<template>
    <div class="flex flex-col min-h-full -mx-4 md:-mx-8 -my-4 md:-my-8">
        <div class="sticky top-0 z-10 shrink-0 bg-base-200 border-b border-base-300">
            <div class="max-w-4xl mx-auto px-6 pt-3">
                <p class="text-xs font-black uppercase tracking-widest text-base-content/40 mb-2">
                    {{ isEditMode ? 'Edit Template' : 'New Template' }}
                </p>
                <div role="tablist" class="tabs tabs-lift tabs-lg">
                    <a v-for="(tab, _index) in tabs" :key="tab.id" role="tab" class="tab"
                        :class="{ 'tab-active': activeTab === tab.id }" @click="setActiveTab(tab.id)">
                        {{ tab.label }}
                    </a>
                </div>
            </div>
        </div>

        <!-- Tab content -->
        <div class="flex-grow mt-5">
            <div class="max-w-4xl mx-auto p-6">
                <div class="grid grid-cols-1 gap-4">

                    <!-- DETAILS TAB -->
                    <div v-show="activeTab === 'details'">
                        <div class="card bg-base-100 border border-base-300 shadow-sm">
                            <div class="card-body gap-5">
                                <h2 class="card-title text-sm">
                                    <span class="badge badge-primary">01</span> Template Details
                                </h2>

                                <!-- Contract Kind -->
                                <fieldset class="fieldset p-0 border-none">
                                    <legend class="fieldset-legend">Contract Type</legend>
                                    <div class="grid grid-cols-2 gap-3 mt-1">
                                        <label class="card border-2 cursor-pointer transition-all" :class="form.contract_kind === TemplateType.frameContract
                                            ? 'border-primary bg-primary/5'
                                            : 'border-base-300 hover:border-base-content/20'">
                                            <input type="radio" v-model="form.contract_kind" :value="TemplateType.frameContract"
                                                class="hidden" />
                                            <div class="card-body p-4 gap-1">
                                                <span class="card-title text-sm">Frame Contract</span>
                                                <p class="text-xs text-base-content/60 font-normal">Top-level agreement
                                                    that groups subcontracts</p>
                                            </div>
                                        </label>
                                        <label class="card border-2 cursor-pointer transition-all" :class="form.contract_kind === TemplateType.subContract
                                            ? 'border-primary bg-primary/5'
                                            : 'border-base-300 hover:border-base-content/20'">
                                            <input type="radio" v-model="form.contract_kind" :value="TemplateType.subContract"
                                                class="hidden" />
                                            <div class="card-body p-4 gap-1">
                                                <span class="card-title text-sm">Subcontract</span>
                                                <p class="text-xs text-base-content/60 font-normal">Scoped agreement
                                                    under a frame contract</p>
                                            </div>
                                        </label>
                                    </div>
                                </fieldset>

                                <fieldset class="fieldset p-0 border-none">
                                    <legend class="fieldset-legend">Global Name</legend>
                                    <input v-model="form.name" class="input input-bordered w-full" type="text"
                                        required />
                                </fieldset>

                                <fieldset class="fieldset p-0 border-none">
                                    <legend class="fieldset-legend">Base Description</legend>
                                    <textarea v-model="form.description" class="textarea textarea-bordered w-full h-24"
                                        required></textarea>
                                </fieldset>

                                <!-- Subcontracts (only for frame contracts) -->
                                <fieldset v-if="form.contract_kind === TemplateType.frameContract"
                                    class="fieldset p-0 border-none">
                                    <legend class="fieldset-legend cursor-pointer select-none inline-flex items-center gap-1.5"
                                        @click="showSubcontractPicker = !showSubcontractPicker">
                                        Subcontract Templates
                                        <svg xmlns="http://www.w3.org/2000/svg"
                                            class="w-3 h-3 transition-transform duration-200 opacity-60"
                                            :class="{ 'rotate-180': showSubcontractPicker }"
                                            fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                d="M19 9l-7 7-7-7" />
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
                                    <div v-if="form.subcontract_template_dids.length" class="flex flex-wrap gap-2 mt-3">
                                        <div v-for="did in form.subcontract_template_dids" :key="did"
                                            class="badge badge-primary badge-outline gap-1 py-3">
                                            <span>{{ getSubcontractTemplateName(did) }}</span>
                                            <button type="button" @click="removeSubcontractTemplate(did)"
                                                class="text-error hover:opacity-70 transition-opacity">✕</button>
                                        </div>
                                    </div>
                                    <p v-else class="fieldset-label mt-2">No subcontract templates selected yet.</p>
                                </fieldset>
                            </div>
                        </div>
                    </div>

                    <!-- SEMANTIC RULES TAB -->
                    <div v-show="activeTab === 'semantic'">
                        <div class="card bg-base-100 border border-base-300 shadow-sm">
                            <div class="card-body gap-5">
                                <h2 class="card-title text-sm">
                                    <span class="badge badge-secondary">02</span> Semantic Rules
                                </h2>
                                <SemanticRulesEditor />
                            </div>
                        </div>
                        
                    </div>

                    <!-- CLAUSES TAB -->
                    <div v-show="activeTab === 'clauses'">
                        <div class="card bg-base-100 border border-base-300 shadow-sm">
                            <div class="card-body gap-5">
                                <h2 class="card-title text-sm">
                                    <span class="badge badge-primary">03</span> Clauses
                                </h2>

                                <div class="card bg-base-200 shadow-none border border-dashed border-base-300">
                                    <div class="card-body gap-3 p-4">
                                        <input v-model="newClause.title"
                                            class="input input-sm input-bordered bg-base-100 font-bold w-full"
                                            placeholder="Clause Title" />
                                        <textarea v-model="newClause.description"
                                            class="textarea textarea-sm textarea-bordered bg-base-100 w-full"
                                            placeholder="Legal text content"></textarea>
                                        <div>
                                            <p class="label-text text-xs text-base-content/50 mb-2">Semantic Rules for
                                                this Clause</p>
                                            <div class="flex flex-wrap gap-2">
                                                <p v-if="!semanticConditions.length"
                                                    class="text-xs text-base-content/50 italic">No semantic rules yet.
                                                </p>
                                                <form class="flex flex-wrap gap-2">
                                                    <input v-for="(rule, idx) in semanticConditions" :key="idx"
                                                        class="btn btn-xs" type="checkbox" :aria-label="rule.conditionName"
                                                        :value="rule.conditionId" v-model="selectedConditionIds" />
                                                </form>
                                            </div>
                                        </div>
                                        <div class="card-actions justify-end">
                                            <button @click="addClause" type="button" class="btn btn-primary btn-sm px-8"
                                                :disabled="!newClause.title || !newClause.description">Add
                                                Clause</button>
                                        </div>
                                    </div>
                                </div>

                                <div class="space-y-2">
                                    <div v-for="(clause, index) in form.clauses" :key="index"
                                        class="card border border-base-300 group hover:border-primary transition-colors">
                                        <div class="card-body p-4 gap-2">
                                            <div class="flex justify-between items-start gap-2">
                                                <div class="flex-1">
                                                    <h4
                                                        class="font-black text-sm uppercase tracking-tight text-primary">
                                                        {{ clause.title }}</h4>
                                                    <p
                                                        class="text-xs text-base-content/70 mt-1 leading-relaxed whitespace-pre-wrap">
                                                        {{ clause.description }}</p>
                                                    <div class="flex flex-wrap gap-1 mt-2">
                                                        <span v-for="cid in clause.conditionIds" :key="cid"
                                                            class="text-primary">{{
                                                            getConditionName(cid) }}</span>
                                                    </div>
                                                </div>
                                                <button @click="removeClause(index)"
                                                    class="btn btn-ghost btn-xs text-error opacity-0 group-hover:opacity-100 transition-opacity shrink-0">✕</button>
                                            </div>
                                        </div>
                                    </div>
                                    <div v-if="!form.clauses?.length"
                                        class="text-center py-8  text-xs text-base-content/40 italic">
                                        No clauses defined yet.
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- BUILDER TAB -->
                    <div v-show="activeTab === 'builder'">
                        <div class="card bg-base-100 border border-base-300 shadow-sm">
                            <div class="card-body">
                                <h2 class="card-title text-sm">Builder</h2>
                                <TemplateEditor />
                            </div>
                        </div>
                        <AddBlockModal />
                    </div>

                    <!-- META TAB -->
                    <div v-show="activeTab === 'meta'">
                        <div class="card bg-base-100 border border-base-300 shadow-sm">
                            <div class="card-body">
                                <h2 class="card-title text-sm">Meta Data</h2>
                                <p class="text-sm text-base-content/60">TODO</p>
                            </div>
                        </div>
                    </div>

                </div>
            </div>
        </div>

        <!-- Pinned Footer -->
        <div class="sticky bottom-0 shrink-0 border-t border-base-300 bg-base-100">
            <div class="max-w-4xl mx-auto px-6 py-3 flex flex-col md:flex-row gap-3">
                <button class="btn btn-ghost md:w-32" @click="cancel">Cancel</button>
                <button @click="submit" class="btn btn-primary flex-1" :disabled="isSubmitting">
                    <span v-if="isSubmitting" class="loading loading-spinner loading-sm"></span>
                    {{ isEditMode ? 'Update Template' : 'Publish to Repository' }}
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useContractTemplateController } from './NewContractTemplate.ts'
import { useTemplateEditorUiStore } from '@template-repository/store/templateEditorUiStore.ts'
import TemplateEditor from '@template-repository/components/TemplateEditor.vue'
import AddBlockModal from '@template-repository/components/AddBlockModal.vue'
import SemanticRulesEditor from '@template-repository/components/SemanticRulesEditor.vue'
import { useTemplateDraftStore } from '@template-repository/store/templateDraftStore'
import { storeToRefs } from 'pinia'
import { TemplateType } from '@template-repository/models/contract-templace'

const templateEditorUiStore = useTemplateEditorUiStore()
const templateDraftStore = useTemplateDraftStore()
const { activeTab, tabs } = storeToRefs(templateEditorUiStore)
const { setActiveTab } = templateEditorUiStore
const { semanticConditions } = storeToRefs(templateDraftStore)

function getConditionName(conditionId: string): string {
  const c = semanticConditions.value.find((x) => x.conditionId === conditionId)
  return c?.conditionName ?? conditionId
}

const props = defineProps<{
    did?: string
    document_number?: number
    version?: number
}>()

const showSubcontractPicker = ref(false)

const {
    form,
    isEditMode,
    isSubmitting,
    submit,
    cancel,
    newClause,
    addClause,
    removeClause,
    selectedConditionIds,
    subcontractSearchQuery,
    filteredSubcontractTemplates,
    getSubcontractTemplateName,
    addSubcontractTemplate,
    removeSubcontractTemplate,
} = useContractTemplateController(props.did, props.document_number, props.version)
</script>
