<script setup lang="ts">
import type { ContractTemplateApprovalTask } from '@/models/contract-template-approval-task'
import { useContractTemplatesStore } from '@/stores/contract-templates-store'
import { computed } from 'vue'

const props = defineProps<{
  items: ContractTemplateApprovalTask[]
}>()

const templatesStore = useContractTemplatesStore()
const templates = computed(() => templatesStore.contractTemplates)

function getTemplateName(item: ContractTemplateApprovalTask) {
  return templates.value.find(template => template.did === item.did)?.name ?? 'Nameless Template'
}
</script>

<template>
  <ul class="list">
    <li v-for="item in items" class="list-row">
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
                  name: 'templates.view',
                  params: { did: item.did },
                }"
                class="btn btn-sm btn-primary rounded-box"
              >
                View
              </RouterLink>
              <RouterLink
                :to="{
                  name: 'templates.edit',
                  params: { did: item.did },
                }"
                class="btn btn-sm btn-secondary rounded-box gap-2"
              >
                Edit
              </RouterLink>
              <RouterLink
                v-if="item.state === 'OPEN'"
                :to="{ name: 'templates.approve', params: { did: item.did } }"
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
