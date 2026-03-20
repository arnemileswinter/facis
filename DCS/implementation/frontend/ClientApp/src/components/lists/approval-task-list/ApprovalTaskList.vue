<script setup lang="ts">
import type { ContractTemplateApprovalTask } from '@/models/contract-template-approval-task'
import { ROUTES } from '@/router/router'
import { useContractTemplatesStore } from '@/stores/contract-templates-store'
import { TemplateState } from '@/types/contract-template-state'
import { computed } from 'vue'

const props = defineProps<{
  items: ContractTemplateApprovalTask[]
}>()

const templatesStore = useContractTemplatesStore()

const sortedItems = computed(() =>
  props.items.sort((taskA, taskB) =>
    new Date(taskA.created_at).getTime() < new Date(taskB.created_at).getTime() ? 1 : -1,
  ),
)

const getTemplateName = (item: ContractTemplateApprovalTask) => {
  return templatesStore.contractTemplates.find((template) => template.did === item.did)?.name ?? 'Nameless Template'
}

const getTemplateState = (item: ContractTemplateApprovalTask) => {
  return templatesStore.contractTemplates.find((template) => template.did === item.did)?.state
}
</script>

<template>
  <ul class="list">
    <li v-for="item in sortedItems" class="list-row">
      <div class="list-col-grow card bg-base-200 card-border hover:bg-base-300">
        <div class="card-body">
          <h2 class="card-title justify-between">
            <div>Approval Task for Contract Template: {{ getTemplateName(item) }}</div>
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
                v-if="item.state === 'OPEN' && getTemplateState(item) === TemplateState.reviewed"
                :to="{ name: ROUTES.TEMPLATES.APPROVE, params: { did: item.did } }"
                class="btn btn-sm btn-primary rounded-box gap-2"
              >
                Approve
              </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </li>
  </ul>
</template>
