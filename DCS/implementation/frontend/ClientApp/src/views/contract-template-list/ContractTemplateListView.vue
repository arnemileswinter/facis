<template>
  <div class="flex justify-between mb-8">
    <h2 class="text-2xl/7 font-bold sm:truncate sm:text-3xl sm:tracking-tight">
      {{ $route.meta.name }}
    </h2>

    <div role="tablist" class="tabs tabs-box tabs-sm">
      <template v-for="(tab, index) in tabs" :key="index">
      <button
        v-show="isTabVisible(tab, index)"
        role="tab"
        class="tab"
        :class="{ 'tab-active': activeTab === index, 'tab-disabled': isTabDisabled(tab) }"
        :disabled="isTabDisabled(tab)"
        @click="activeTab = index"
      >
        {{ tab }}
      </button></template>
    </div>

    <RouterLink
      :to="{ name: 'templates.new' }"
      class="btn rounded-box self-end btn-secondary gap-2"
      #default="{ route }"
    >
      {{ route.meta.name }}
    </RouterLink>
  </div>
  <div class="p-6 bg-base-100">
    <div v-if="loading">Lade Templates...</div>
    <div v-else-if="error">{{ error }}</div>
    <div v-else>
      <div v-if="activeTab === 0">
        <TemplateList :items="templates"
        :has-review-task="hasReviewTask"
        :has-approval-task="hasApprovalTask"
      />
      </div>
      <div v-if="activeTab === 1">
        <ReviewTaskList :items="reviewTasks" />
      </div>
      <div v-if="activeTab === 2">
        <ApprovalTaskList :items="approvalTasks" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTemplateTable } from './ContractTemplateListController'
import TemplateList from '../../components/lists/template-list/TemplateList.vue'
import { ref } from 'vue'
import ApprovalTaskList from '@/components/lists/approval-task-list/ApprovalTaskList.vue'
import ReviewTaskList from '@/components/lists/review-task-list/ReviewTaskList.vue'
import type { UserRole } from '@/types/user-role'

const tabs = ['Contract Templates', 'Review Tasks', 'Approval Tasks']
const activeTab = ref(0)
const tabRoles: Record<string, UserRole> = {'Review Tasks': 'TEMPLATE_REVIEWER', 'Approval Tasks': 'TEMPLATE_APPROVER'}

const { templates, reviewTasks, approvalTasks, roles, loading, error, hasReviewTask, hasApprovalTask } = useTemplateTable()

function isTabVisible(tab: string, index: number) {
  return (
    (index === 0) ||
    (tab === 'Review Tasks' && roles.value?.includes(tabRoles[tab]!)) ||
    (tab === 'Approval Tasks' && roles.value?.includes(tabRoles[tab]!))
  )
}

function isTabDisabled(tab: string) {
  return (
    (tab === 'Review Tasks' && reviewTasks.value.length < 1) ||
    (tab === 'Approval Tasks' && approvalTasks.value.length < 1)
  )
}
</script>
