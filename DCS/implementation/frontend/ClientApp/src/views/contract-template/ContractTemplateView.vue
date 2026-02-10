<template>
  <div class="flex justify-between
">
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
      <DataTable :items="templates" :headers="headers">
        <template #extraCols="{ item }">
          <ContractTemplateTableRowData :template="item" />
        </template>
      </DataTable>
    </div>
  </div>
</template>

<script setup lang="ts">
import ContractTemplateTableRowData from '../../components/tables/ContractTemplateTableRowData.vue'
import DataTable from '../../components/tables/DataTable.vue'
import { useTemplateTable } from './ContractTemplateController'

const { templates, loading, error, refresh } = useTemplateTable()
const headers = [
  'did',
  'documentNumber',
  'name',
  'state',
  'version',
  'createdBy',
  'createdAt',
  'description',
  'meta_data',
]
</script>
