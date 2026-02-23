<script setup lang="ts">
import type { ContractTemplate } from '../../../models/contract-template'

const props = defineProps<{
  item: ContractTemplate
}>()
</script>

<template>
  <li class="list-row">
    <div class="list-col-grow card bg-base-200 card-border hover:bg-base-300">
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
        <div class="card-actions justify-between">
          <div class="card-actions justify-start">
            <input type="checkbox" class="checkbox checkbox-xl checkbox-accent" />
          </div>
          <div class="card-actions justify-end">
            <button class="btn btn-outline rounded-box btn-info">View</button>
            <RouterLink
              :to="{
                name: 'templates.edit',
                params: { did: item.did },
                query: { document_number: item.document_number, version: item.version },
              }"
              class="btn btn-outline rounded-box btn-warning gap-2"
            >
              Edit
            </RouterLink>
          </div>
        </div>
      </div>
    </div>
  </li>
</template>
