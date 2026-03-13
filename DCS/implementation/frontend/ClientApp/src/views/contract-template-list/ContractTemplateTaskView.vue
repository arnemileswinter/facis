<script setup lang="ts">
import ApprovalTaskList from '@/components/lists/approval-task-list/ApprovalTaskList.vue'
import ReviewTaskList from '@/components/lists/review-task-list/ReviewTaskList.vue'
import { useTemplateTable } from './ContractTemplateListController'

const { reviewTasks, approvalTasks } = useTemplateTable()
</script>

<template>
  <h2 class="text-2xl/7 font-bold sm:truncate sm:text-3xl sm:tracking-tight">
    {{ $route.meta.name }}
  </h2>

  <div v-if="$route.name === 'templates.tasks.review'">
    <ReviewTaskList :items="reviewTasks" />
  </div>
  <div v-else-if="$route.name === 'templates.tasks.approval'">
    <ApprovalTaskList :items="approvalTasks" />
  </div>
</template>
