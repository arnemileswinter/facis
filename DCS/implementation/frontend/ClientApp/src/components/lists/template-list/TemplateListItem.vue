<script setup lang="ts">
import type { PartialContractTemplate } from '@/models/contract-template'
import { useAuthStore } from '@/stores/auth-store'
import { TemplateState } from '@/types/contract-template-state'
import { toProperCase } from '@/utils/string'
import { computed } from 'vue'

const props = defineProps<{
  item: PartialContractTemplate
  hasReviewTask: boolean
  hasApprovalTask: boolean
}>()

const authStore = useAuthStore()

const canEdit = computed(() => {
  return (
    ((props.item.state === TemplateState.draft || props.item.state === TemplateState.rejected) &&
      props.item.created_by === authStore.user?.username) ||
    (props.item.state === TemplateState.submitted && props.hasReviewTask)
  )
})

</script>

<template>
  <li class="list-row min-w-0 w-full">
    <div class="list-col-grow card bg-base-200 card-border hover:bg-base-300 min-w-0 w-full">
      <div class="card-body min-w-0">
        <h2 class="card-title justify-between">
          <div class="flex gap-8 h-full">
            <div>Name: {{ item.name }}</div>
            <div class="badge badge-md badge-accent h-full">{{ toProperCase(item.template_type) }}</div>
          </div>
          <div class="badge badge-secondary">{{ item.state }}</div>
        </h2>
        <div class="flex justify-between">
          <div v-if="item.document_number">Document number: {{ item.document_number }}</div>
          <div v-if="item.version">Version: {{ item.version }}</div>
        </div>
        <div class="flex justify-between min-w-0">
          <div>Creation date: {{ new Date(item.created_at).toLocaleDateString() }}</div>
          <div v-if="item.description" class="px-10 flex-1 min-w-0 truncate">
            {{ item.description }}
          </div>
          <div class="card-actions justify-end">
            <RouterLink
              :to="{ name: 'templates.view', params: { did: item.did } }"
              class="btn btn-sm btn-primary rounded-box"
            >
              View
            </RouterLink>
            <RouterLink
              v-if="canEdit"
              :to="{
                name: 'templates.edit',
                params: { did: item.did },
              }"
              class="btn btn-sm btn-primary rounded-box gap-2"
            >
              Edit
            </RouterLink>
            <RouterLink
              v-if="item.state === TemplateState.submitted && hasReviewTask"
              :to="{ name: 'templates.review', params: { did: item.did } }"
              class="btn btn-sm btn-primary rounded-box gap-2"
            >
              Review
            </RouterLink>
            <RouterLink
              v-if="item.state === TemplateState.reviewed && hasApprovalTask"
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
</template>
