<script setup lang="ts">
import type { ContractTemplateReviewTask } from '@/models/contract-template-review-task'
import { ROUTES } from '@/router/router'
import { useAuthStore } from '@/stores/auth-store'
import { useContractTemplatesStore } from '@/stores/contract-templates-store'
import { TemplateState } from '@/types/contract-template-state'
import { computed } from 'vue'

const props = defineProps<{
  items: ContractTemplateReviewTask[]
}>()

const templatesStore = useContractTemplatesStore()
const authStore = useAuthStore()

const sortedItems = computed(() =>
  props.items.sort((taskA, taskB) =>
    new Date(taskA.created_at).getTime() < new Date(taskB.created_at).getTime() ? 1 : -1,
  ),
)

const getTemplateName = (item: ContractTemplateReviewTask) => {
  return templatesStore.contractTemplates.find((template) => template.did === item.did)?.name ?? 'Nameless Template'
}

const canEdit = (item: ContractTemplateReviewTask) => {
  const template = templatesStore.contractTemplates.find((template) => template.did === item.did)
  const state = template?.state
  return (
    (template?.created_by === authStore.user?.username &&
      (state === TemplateState.draft || state === TemplateState.rejected)) ||
    state === TemplateState.submitted
  )
}
</script>

<template>
  <ul class="list">
    <li v-for="item in sortedItems" class="list-row">
      <div class="list-col-grow card bg-base-200 card-border hover:bg-base-300">
        <div class="card-body">
          <h2 class="card-title justify-between">
            <div>Review Task for Contract Template: {{ getTemplateName(item) }}</div>
            <div class="badge badge-secondary">{{ item.state }}</div>
          </h2>
          <div class="flex justify-between">
            <div v-if="item.document_number">Document number: {{ item.document_number }}</div>
            <div v-if="item.version">Version: {{ item.version }}</div>
          </div>
          <div class="flex justify-between">
            <div>Creation date: {{ new Date(item.created_at).toLocaleDateString() }}</div>
            <div class="card-actions justify-end">
              <RouterLink
                :to="{
                  name: ROUTES.TEMPLATES.VIEW,
                  params: { did: item.did },
                }"
                class="btn btn-sm btn-primary rounded-box"
              >
                View
              </RouterLink>
              <RouterLink
                v-if="canEdit(item)"
                :to="{
                  name: ROUTES.TEMPLATES.EDIT,
                  params: { did: item.did },
                }"
                class="btn btn-sm btn-secondary rounded-box gap-2"
              >
                Edit
              </RouterLink>
              <RouterLink
                v-if="item.state === 'OPEN'"
                :to="{ name: ROUTES.TEMPLATES.REVIEW, params: { did: item.did } }"
                class="btn btn-sm btn-primary rounded-box gap-2"
              >
                Review
              </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </li>
  </ul>
</template>
