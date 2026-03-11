<script setup lang="ts">
import UserSelectionDialog from '@/components/UserSelectionDialog.vue'
import type { ContractTemplateSubmitRequest } from '@/models/requests/template-request'
import type { SelectedUserRole } from '@/models/user'
import { ContractTemplateService } from '@/services/contract-template-service'
import { TemplateState } from '@/types/contract-template-state'
import type { PartialContractTemplate } from '../../../models/contract-template'

const props = defineProps<{
  item: PartialContractTemplate
}>()

async function submitTemplate(result: SelectedUserRole[]) {
  const reviewers = result.filter((user) => user.role === 'TEMPLATE_REVIEWER').map((user) => user.user.id)
  const approver = result.find((user) => user.role === 'TEMPLATE_APPROVER')?.user.id
  const request: ContractTemplateSubmitRequest = {
    did: props.item.did,
    updated_at: props.item.updated_at,
    reviewers: reviewers,
    approver: approver!,
  }
  const response = await ContractTemplateService.submit(request)
  if (response) {
    console.log('Successful submitted.')
  }
}
</script>

<template>
  <li class="list-row">
    <div class="list-col-grow card bg-base-200 card-border hover:bg-base-300">
      <div class="card-body">
        <h2 class="card-title justify-between">
          <div>Name: {{ item.name }}</div>
          <div class="badge badge-secondary">{{ item.state }}</div>
        </h2>
        <div class="flex justify-between">
          <div>Document number: {{ item.document_number }}</div>
          <div>Version: {{ item.version }}</div>
        </div>
        <div class="flex justify-between">
          <div>Creation date: {{ new Date(item.created_at).toLocaleDateString() }}</div>
          <div class="card-actions justify-end">
            <RouterLink
              :to="{ name: 'templates.view', params: { did: item.did } }"
              class="btn btn-sm btn-primary rounded-box"
            >
              View
            </RouterLink>
            <UserSelectionDialog
              v-if="item.state === TemplateState.draft"
              @submit="submitTemplate"
              class="btn btn-sm btn-secondary rounded-box"
            />
            <RouterLink
              v-if="item.state === TemplateState.draft || item.state === TemplateState.rejected"
              :to="{
                name: 'templates.edit',
                params: { did: item.did },
              }"
              class="btn btn-sm btn-primary rounded-box gap-2"
            >
              Edit
            </RouterLink>
            <RouterLink
              v-if="item.state === TemplateState.submitted"
              :to="{ name: 'templates.review', params: { did: item.did } }"
              class="btn btn-sm btn-primary rounded-box gap-2"
            >
              Review
            </RouterLink>
            <RouterLink
              v-if="item.state === TemplateState.reviewed"
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
