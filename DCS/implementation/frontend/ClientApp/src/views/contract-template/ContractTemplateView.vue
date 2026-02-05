<template>
  <div>
    <div v-if="loading">Lade Templates...</div>
    <div v-else-if="error">{{ error }}</div>
    <div v-else>
      <DataTable :items="templates" />
    </div>
  </div>
</template>

<script setup lang="ts">
import DataTable from '../../components/DataTable.vue';
import { useTemplateTable } from './ContractTemplateController';

const { templates, loading, error, refresh } = useTemplateTable()
</script>
