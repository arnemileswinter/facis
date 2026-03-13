<script setup lang="ts">
import SubmitContractTemplateUserSelectionDialog from '@/components/SubmitContractTemplateUserSelectionDialog.vue'
import type { ContractTemplateSubmitRequest } from '@/models/requests/template-request'
import type { SelectedUserRole } from '@/models/user'
import { ContractTemplateService } from '@/services/contract-template-service'
import { TemplateState } from '@/types/contract-template-state'
import type { PartialContractTemplate } from '../../../models/contract-template'
import { useRouter } from 'vue-router'

const props = defineProps<{
  item: PartialContractTemplate
  hasReviewTask: boolean
  hasApprovalTask: boolean
}>()

const router = useRouter()

async function submitTemplate(result: SelectedUserRole[]) {
  const reviewers = result.filter((user) => user.role === 'TEMPLATE_REVIEWER').map((user) => user.user.username)
  const approver = result.find((user) => user.role === 'TEMPLATE_APPROVER')?.user.username
  const request: ContractTemplateSubmitRequest = {
    did: props.item.did,
    updated_at: props.item.updated_at,
    reviewers: reviewers,
    approver: approver!,
  }
  const response = await ContractTemplateService.submit(request)
  if (response) {
    router.go(0)
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
          <div v-if="item.document_number">Document number: {{ item.document_number }}</div>
          <div v-if="item.version">Version: {{ item.version }}</div>
        </div>
        <div class="flex justify-between">
          <div>Creation date: {{ new Date(item.created_at).toLocaleDateString() }}</div>
          <div v-if="item.description" class="px-4 max-w-1/12 whitespace-nowrap overflow-hidden text-ellipsis">
            {{ item.description }}
          </div>
          <div class="flex-1"></div>
          <div class="card-actions justify-end">
            <RouterLink
              :to="{ name: 'templates.view', params: { did: item.did } }"
              class="btn btn-sm btn-primary rounded-box"
            >
              View
            </RouterLink>
            <SubmitContractTemplateUserSelectionDialog
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
