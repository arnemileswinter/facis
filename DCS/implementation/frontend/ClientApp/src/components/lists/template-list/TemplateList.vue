<script setup lang="ts">
import { computed, ref } from 'vue'
import type { ContractTemplate } from '../../../models/contract-template'
import ListSort from '../ListSort.vue'
import TemplateListItem from './TemplateListItem.vue'
import ListSearch from '../ListSearch.vue';

const props = defineProps<{
  items: ContractTemplate[]
}>()

const sorter = new Map([
  ['created_at', 'Creation date'],
  ['state', 'Status'],
])
const defaultSort = sorter.keys().next().value!
const sortBy = ref(defaultSort)
const sortOrder = ref(1)

const itemsSorted = computed(() => {
  if (!sorter.has(sortBy.value)) {
    return props.items
  }
  return props.items.slice().sort((a, b) => {
    let aSortValue = a[sortBy.value as keyof ContractTemplate]
    let bSortValue = b[sortBy.value as keyof ContractTemplate]
    if (sortBy.value === defaultSort) {
      aSortValue = new Date(aSortValue).getTime()
      bSortValue = new Date(bSortValue).getTime()
    }
    const result = aSortValue > bSortValue ? 1 : -1
    return sortOrder.value * result
  })
})
</script>

<template>
  <ul class="list">
    <li class="tracking-wide px-4 flex justify-between">
      <ListSearch class="grow" />
      <ListSort :sorter="sorter" v-model:sort-by="sortBy" v-model:sort-order="sortOrder" />
    </li>
    <TemplateListItem v-for="item in itemsSorted" :key="item.did" :item="item" />
  </ul>
</template>
