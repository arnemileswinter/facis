<template>
    <div class="max-w-4xl mx-auto p-6 space-y-8">
        <div class="flex items-center justify-between border-b border-base-300 pb-6">
            <div>
                <h1 class="text-3xl font-black uppercase tracking-tighter">
                    {{ isEditMode ? 'Edit Template' : 'New Template' }}
                </h1>
            </div>
        </div>

        <div class="grid grid-cols-1 gap-6">
            <details class="collapse border border-base-300 bg-base-100 shadow-sm" name="my-accordion" open>
                <summary class="collapse-title font-bold text-lg flex items-center gap-2">
                    <span class="text-primary">01</span> Template Details
                </summary>
                <div class="collapse-content space-y-4">
                    <div class="form-control w-full">
                        <label class="label"><span class="label-text-alt uppercase font-bold opacity-50">Global
                                Name</span></label>
                        <input v-model="form.name" class="input input-bordered validator w-full" type="text" required />
                    </div>
                    <div class="form-control w-full">
                        <label class="label"><span class="label-text-alt uppercase font-bold opacity-50">Base
                                Description</span></label>
                        <textarea v-model="form.description" class="textarea textarea-bordered w-full h-24"
                            required></textarea>
                    </div>
                </div>
            </details>

            <details class="collapse border border-base-300 bg-base-100 shadow-sm" name="my-accordion" open>
                <summary class="collapse-title font-bold text-lg flex items-center gap-2">
                    <span class="text-primary">02</span> Clauses
                </summary>
                <div class="collapse-content space-y-6">
                    <div class="flex flex-col gap-3 p-4 bg-base-200/50 rounded-xl border border-dashed border-base-300">
                        <input v-model="newClause.title" class="input input-sm input-bordered w-full font-bold"
                            placeholder="Clause Title " />
                        <textarea v-model="newClause.content" class="textarea textarea-sm textarea-bordered w-full"
                            placeholder="Legal text content "></textarea>
                        <button @click="addClause" type="button" class="btn btn-primary btn-sm self-end px-8"
                            :disabled="!newClause.title || !newClause.content">
                            Add Clause
                        </button>
                    </div>

                    <div class="space-y-3">
                        <div v-for="(clause, index) in form.clauses" :key="index"
                            class="flex flex-col p-4 border border-base-300 rounded-xl bg-base-50 relative group transition-all hover:border-primary/50">
                            <div class="flex justify-between items-start">
                                <div>
                                    <h4 class="font-black text-sm uppercase tracking-tight text-primary">{{ clause.title
                                        }}</h4>
                                    <p class="text-xs opacity-80 mt-2 leading-relaxed whitespace-pre-wrap">{{
                                        clause.content }}</p>
                                </div>
                                <button @click="removeClause(index)"
                                    class="btn btn-ghost btn-xs text-error opacity-0 group-hover:opacity-100 transition-opacity">✕</button>
                            </div>
                        </div>
                        <div v-if="!form.clauses?.length"
                            class="text-center py-8 border-2 border-dashed border-base-200 rounded-xl text-xs opacity-40 italic">
                            No clauses defined yet.
                        </div>
                    </div>
                </div>
            </details>

            <details class="collapse border border-secondary/30 bg-base-100 shadow-sm overflow-visible"
                name="my-accordion" open>
                <summary class="collapse-title font-bold text-lg flex items-center gap-2">
                    <span class="text-secondary">03</span> Semantic Rules
                </summary>
                <div class="collapse-content space-y-6">
                    <div v-if="suggestions.length">
                        <label class="text-[10px] uppercase font-black tracking-widest opacity-50 mb-3 block">Suggested
                            by Semantic Hub</label>
                        <div class="flex flex-wrap gap-2">
                            <button v-for="sug in suggestions" :key="sug.attr" @click="addRuleFromSuggestion(sug)"
                                class="btn btn-outline btn-secondary btn-xs normal-case hover:bg-secondary hover:text-white transition-all">
                                + {{ sug.attr }}
                            </button>
                        </div>
                    </div>

                    <div
                        class="grid grid-cols-1 md:grid-cols-12 gap-2 p-4 bg-secondary/5 rounded-xl border border-secondary/20 items-end">
                        <div class="md:col-span-4">
                            <label class="label-text text-[10px] uppercase font-bold ml-1 opacity-60">Label</label>
                            <input v-model="newRule.attr" type="text"
                                class="input input-bordered input-sm w-full mt-1" />
                        </div>
                        <div class="md:col-span-3">
                            <label class="label-text text-[10px] uppercase font-bold ml-1 opacity-60">Type</label>
                            <select v-model="newRule.op" class="select select-bordered select-sm w-full mt-1">
                                <option value="MIN">Date</option>
                                <option value="MAX">Text </option>
                                <option value="EQUALS">Decimal</option>
                            </select>
                        </div>
                        <div class="md:col-span-3">
                            <select v-model="newRule.op" class="select select-bordered select-sm w-full mt-1">
                                <option value="MIN">Date</option>
                                <option value="MAX">Text </option>
                                <option value="EQUALS">Decimal</option>
                            </select>
                        </div>
                        <select v-model="form.semantic_rules" class="select select-bordered select-sm w-full mt-1">
                            <option value="MIN">required</option>
                            <option value="MAX">optional</option>
                        </select>
                        <div class="md:col-span-2">
                            <button @click="addNewCustomRule" class="btn btn-secondary btn-sm w-full mt-1"
                                :disabled="!newRule.attr">Add</button>
                        </div>
                    </div>

                    <div class="space-y-2">
                        <div v-for="(rule, index) in form.semantic_rules" :key="index"
                            class="flex items-center gap-4 p-3 bg-base-100 border border-base-300 border-l-4 border-l-secondary rounded-r-lg group hover:shadow-md transition-all">
                            <div class="bg-secondary/10 p-2 rounded text-secondary">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24"
                                    stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3"
                                        d="M13 10V3L4 14h7v7l9-11h-7z" />
                                </svg>
                            </div>
                            <div class="flex-1">
                                <div class="flex items-center gap-2">
                                    <span class="text-sm font-black font-mono uppercase">{{ rule.attr }}</span>
                                    <span class="badge badge-ghost badge-xs">{{ rule.op }}</span>
                                </div>
                                <p class="text-[10px] font-bold opacity-50">Target: {{ rule.val }}</p>
                            </div>
                            <button @click="removeRule(index)"
                                class="btn btn-ghost btn-xs text-error opacity-0 group-hover:opacity-100 transition-opacity">✕</button>
                        </div>
                    </div>
                </div>
            </details>

            <div class="flex flex-col md:flex-row gap-3 pt-6 border-t border-base-300">
                <button class="btn btn-ghost md:w-32" @click="cancel">Cancel</button>
                <button @click="submit" class="btn btn-primary flex-1 shadow-lg" :disabled="isSubmitting">
                    <span v-if="isSubmitting" class="loading loading-spinner"></span>
                    {{ isEditMode ? 'Update Template Architecture' : 'Publish Template to Repository' }}
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useContractTemplateController } from './NewContractTemplate.ts'

const {
    form,
    isEditMode,
    isSubmitting,
    submit,
    cancel,
    newClause,
    addClause,
    removeClause,
    newRule,
    suggestions,
    addRuleFromSuggestion,
    addNewCustomRule,
    removeRule
} = useContractTemplateController()
</script>