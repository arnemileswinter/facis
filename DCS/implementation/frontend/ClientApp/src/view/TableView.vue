<script setup lang="ts">
import Table from '../components/Table.vue'
import type { TableItem } from '../models/table-item'

const tableItems: TableItem[] = [{ id: 12, name: 'Tester' }]
</script>

<template>
  <Table :items="tableItems" />
</template>
