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
                                <DetailsEditor ref="detailsEditorRef" />
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
                                <ClausesEditor />
                            </div>
                        </div>
                    </div>

                    <!-- BUILDER TAB -->
                    <div v-show="activeTab === 'builder'">
                        <div class="card bg-base-100 border border-base-300 shadow-sm">
                            <div class="card-body">
                                <h2 class="card-title text-sm">Builder</h2>
                                <BuilderEditor />
                            </div>
                        </div>
                        <AddBlockModal />
                    </div>

                    <!-- META TAB -->
                    <div v-show="activeTab === 'meta'">
                        <div class="card bg-base-100 border border-base-300 shadow-sm">
                            <div class="card-body">
                                <h2 class="card-title text-sm">Meta Data</h2>
                                <MetaDataEditor />
                            </div>
                        </div>
                    </div>

                </div>
            </div>
        </div>

        <!-- Pinned Footer -->
        <div class="sticky bottom-0 shrink-0 border-t border-base-300 bg-base-100">
            <div class="max-w-4xl mx-auto px-6 py-3 flex flex-col md:flex-row gap-3">
                <button class="btn btn-ghost md:w-32" @click="router.back()">Cancel</button>
                <button @click="submit" class="btn btn-primary flex-1" :disabled="isSubmitting">
                    <span v-if="isSubmitting" class="loading loading-spinner loading-sm"></span>
                    {{ isEditMode ? 'Update Template' : 'Publish to Repository' }}
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useTemplateEditorUiStore } from '@template-repository/store/templateEditorUiStore.ts'
import BuilderEditor from '@template-repository/components/BuilderEditor.vue'
import AddBlockModal from '@template-repository/components/builder-editor/AddBlockModal.vue'
import SemanticRulesEditor from '@template-repository/components/SemanticRulesEditor.vue'
import ClausesEditor from '@template-repository/components/ClausesEditor.vue'
import DetailsEditor from '@template-repository/components/DetailsEditor.vue'
import MetaDataEditor from '@template-repository/components/MetaDataEditor.vue'
import { storeToRefs } from 'pinia'

const router = useRouter()
const route = useRoute()

const templateEditorUiStore = useTemplateEditorUiStore()
const { activeTab, tabs } = storeToRefs(templateEditorUiStore)
const { setActiveTab } = templateEditorUiStore

const isEditMode = computed(() => !!route.params.did)
const isSubmitting = ref(false)

const detailsEditorRef = ref<InstanceType<typeof DetailsEditor> | null>(null)

const submit = async () => {
    isSubmitting.value = true
    try {
        const formData = detailsEditorRef.value?.getFormData()
        console.log('Publishing Template to Repository...', formData)
        await new Promise(resolve => setTimeout(resolve, 1500))
        router.push({ name: 'templates.list' })
    } catch (error) {
        console.error('Submission failed', error)
    } finally {
        isSubmitting.value = false
    }
}
</script>
