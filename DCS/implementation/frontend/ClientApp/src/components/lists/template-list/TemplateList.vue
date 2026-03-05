<script setup lang="ts">
import type { PartialContractTemplate } from '@/models/contract-template'
import { useContractTemplateStateFilterStore } from '@/stores/contract-template-state-filter-store'
import { storeToRefs } from 'pinia'
import { computed, ref } from 'vue'
import ListSearch from '../ListSearch.vue'
import ListSort from '../ListSort.vue'
import TemplateListItem from './TemplateListItem.vue'

const props = defineProps<{
  items: PartialContractTemplate[]
}>()

const sorter = new Map([
  ['created_at', 'Creation date'],
  ['state', 'Status'],
])
const defaultSort = sorter.keys().next().value!
const sortBy = ref(defaultSort)
const sortOrder = ref(1)

const stateFilterStore = useContractTemplateStateFilterStore()
const { stateFilters } = storeToRefs(stateFilterStore)

const sortedItems = computed(() => {
  if (!sorter.has(sortBy.value)) {
    return props.items
  }
  return props.items.slice().sort((a, b) => {
    let aSortValue = a[sortBy.value as keyof PartialContractTemplate]
    let bSortValue = b[sortBy.value as keyof PartialContractTemplate]
    if (sortBy.value === defaultSort && sortBy.value === 'created_at' 
      && typeof aSortValue === 'string' && typeof bSortValue === 'string') {
      aSortValue = new Date(aSortValue).getTime()
      bSortValue = new Date(bSortValue).getTime()
    }
    if (aSortValue === undefined || bSortValue === undefined) return sortOrder.value * 1
    const result = aSortValue > bSortValue ? 1 : -1
    return sortOrder.value * result
  })
})

const filteredItems = computed(() => {
  const filters = stateFilters.value
  if (filters.size > 0) {
    return sortedItems.value.filter((item) => filters.has(item.state))
  }
  return sortedItems.value
})
</script>

<template>
  <ul class="list">
    <li class="tracking-wide px-4 flex justify-between">
      <ListSearch class="grow" />
      <ListSort :sorter="sorter" v-model:sort-by="sortBy" v-model:sort-order="sortOrder" />
    </li>
    <TemplateListItem
      v-for="item in filteredItems"
      :key="`${item.did},${item.document_number},${item.version}`"
      :item="item"
    />
  </ul>
</template>
