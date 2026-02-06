<script setup lang="ts">
import type { ContractTemplate } from '../../../models/contract-template'

defineProps<{
  item: ContractTemplate
}>()
</script>

<template>
  <li class="list-row">
    <div class="list-col-grow card bg-base-100 card-border hover:bg-base-300">
      <div class="card-body">
        <h2 class="card-title justify-between">
          <div>Name: {{ item.name }}</div>
          <div class="badge badge-secondary">{{ item.state }}</div>
        </h2>
        <div class="mt-4 flex justify-between">
          <div>Document number: {{ item.document_number }}</div>
          <div>Version: {{ item.version }}</div>
        </div>
        <div class="my-4 flex justify-between">
          <div>Created by: {{ item.created_by }}</div>
          <div>Creation date: {{ new Date(item.created_at).toLocaleDateString() }}</div>
        </div>
        <div class="card-actions justify-end">
          <button class="btn btn-outline btn-info">View</button>
          <button class="btn btn-outline btn-warning">Edit</button>
        </div>
      </div>
    </div>
  </li>
</template>
