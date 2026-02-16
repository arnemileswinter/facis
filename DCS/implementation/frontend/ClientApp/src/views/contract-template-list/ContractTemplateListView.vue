<template>
  <div class="flex justify-between mb-8">
    <h2 class="text-2xl/7 font-bold  sm:truncate sm:text-3xl sm:tracking-tight ">
      Contract Templates
    </h2>

    <RouterLink to="/templates/new" class="btn rounded-box btn self-end btn-secondary gap-2">
      Neues Template
    </RouterLink>
  </div>
  <div>
    <div v-if="loading">Lade Templates...</div>
    <div v-else-if="error">{{ error }}</div>
    <div v-else>
      <TemplateList :items="templates" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTemplateTable } from './ContractTemplateListController'
import TemplateList from '../../components/lists/template-list/TemplateList.vue'

const { templates, loading, error, refresh } = useTemplateTable()
</script>
