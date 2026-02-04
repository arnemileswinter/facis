<script setup lang="ts" generic="T extends TableItem">
import type { TableItem } from '../models/table-item'

defineProps<{
  item: T
}>()
</script>

<template>
  <tr>
    <th class="pl-8">{{ item.id }}</th>
    <td class="pl-8">{{ item.name }}</td>
    <slot name="extraRows" :item="item"></slot>
    <td class="pl-8">
      <button class="cursor-pointer">Edit</button>
    </td>
  </tr>
</template>
