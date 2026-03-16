<script setup lang="ts">
import SubmitContractTemplateUserSelectionDialog from '@/components/SubmitContractTemplateUserSelectionDialog.vue'
import type { PartialContractTemplate } from '@/models/contract-template'
import type { ContractTemplateSubmitRequest } from '@/models/requests/template-request'
import type { SelectedUserRole } from '@/models/user'
import { ContractTemplateService } from '@/services/contract-template-service'
import { useAuthStore } from '@/stores/auth-store'
import { TemplateState } from '@/types/contract-template-state'
import { toProperCase } from '@/utils/string'
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps<{
  item: PartialContractTemplate
  hasReviewTask: boolean
  hasApprovalTask: boolean
}>()

const canEdit = ref(false)

const router = useRouter()
const authStore = useAuthStore()

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

watch(
  canEdit,
  async () => {
    try {
      const template = await ContractTemplateService.retrieveById({ did: props.item.did })
      const creator = template?.created_by
      canEdit.value =
        ((props.item.state === TemplateState.draft || props.item.state === TemplateState.rejected) &&
          authStore.user?.username === creator) ||
        (props.item.state === TemplateState.submitted && props.hasReviewTask)
    } catch (err) {
      console.error('Error:', err)
      canEdit.value = false
    }
  },
  { immediate: true },
)
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
          <div class="flex-none">Creation date: {{ new Date(item.created_at).toLocaleDateString() }}</div>
          <div v-if="item.description" class="px-10 flex-1 min-w-0 truncate">
            {{ item.description }}
          </div>
          <div class="card-actions justify-end flex-none">
            <RouterLink
              :to="{ name: 'templates.view', params: { did: item.did } }"
              class="btn btn-sm btn-primary rounded-box"
            >
              View
            </RouterLink>
            <SubmitContractTemplateUserSelectionDialog
              v-if="item.state === TemplateState.draft || item.state === TemplateState.rejected"
              @submit="submitTemplate"
              class="btn btn-sm btn-secondary rounded-box"
            />
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
