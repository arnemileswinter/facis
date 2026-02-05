<script setup lang="ts" generic="T extends TableItem">
import type { TableItem } from '../models/table-item'

defineProps<{
  item: T
  hideDefault?: boolean
}>()
</script>

<template>
  <tr>
    <template v-if="!hideDefault">
      <th class="pl-8">{{ item.id }}</th>
      <td class="pl-8">{{ item.name }}</td>
    </template>
    <slot name="extraCols" :item="item"></slot>
    <td class="pl-8">
      <button class="cursor-pointer border p-1 hover:bg-gray-200">Edit</button>
    </td>
  </tr>
</template>
