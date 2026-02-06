<template>
  <h2 class="text-2xl/7 font-bold text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight dark:text-white">
    Contract Templates
  </h2>
  <button class="btn">Neu</button>

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
import ContractTemplateTableRowData from '../../components/ContractTemplateTableRowData.vue'
import DataTable from '../../components/DataTable.vue'
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
