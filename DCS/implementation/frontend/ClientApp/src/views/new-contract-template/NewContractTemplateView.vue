<template>
    <div class="max-w-4xl mx-auto p-6">
        <div class="flex items-center justify-between mb-8">
            <div>
                <h1 class="text-3xl font-bold  uppercase tracking-wider">
                    {{ isEditMode ? 'Edit Template' : 'New Template' }}
                </h1>
            </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-1 gap-8">
            <div class="lg:col-span-1 space-y-6">
                <div class="card  border border-white/5 shadow-2xl">
                    <div class="card-body">
                        <div class="form-control">
                            <label class="label"><span class="label-text text-gray-400">Template Name</span></label>
                            <input v-model="form.name" type="text"
                                class="input input-bordered bg-gray-800 border-gray-700 " />
                        </div>

                        <div class="form-control mt-4">
                            <label class="label"><span class="label-text text-gray-400">Description</span></label>
                            <textarea v-model="form.description"
                                class="textarea textarea-bordered bg-gray-800 border-gray-700  h-32"></textarea>

                        </div>

                    </div>
                </div>

                <div class="space-y-6">
                    <button class="btn btn-ghost" @click="cancel">Cancel</button>
                    <button @click="submit" class="btn btn-primary w-full" :disabled="isSubmitting">
                        <span v-if="isSubmitting" class="loading loading-spinner"></span>
                        {{ isEditMode ? 'Update Template' : 'Create Template' }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useContractTemplateController } from './NewContractTemplate.ts'
const { form, isEditMode, isSubmitting, submit, cancel } = useContractTemplateController()
</script>