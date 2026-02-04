<script setup lang="ts" generic="T extends TableItem">
import type { TableItem } from '../models/table-item'
import TableRow from './TableRow.vue'

const props = defineProps<{
  items: T[]
}>()

const headers = ['id', 'name']
</script>

<template>
  <div class="overflow-x-auto">
    <table class="m-4">
      <thead>
        <tr>
          <template v-for="header in headers">
            <th>
              <button>
                <span>{{ header.toLocaleUpperCase('de') }}</span>
              </button>
            </th>
          </template>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <template v-for="item in items">
          <TableRow :item="item" />
        </template>
      </tbody>
    </table>
  </div>
  <div></div>
</template>
