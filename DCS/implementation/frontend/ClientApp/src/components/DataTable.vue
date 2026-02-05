<script setup lang="ts" generic="T extends TableItem">
import { computed, ref } from 'vue'
import type { TableItem } from '../models/table-item'
import Pagination from './Pagination.vue'
import TableRow from './TableRow.vue'

const { items, headers = [] } = defineProps<{
  items: T[]
  readonly headers?: string[]
}>()

const itemsPerPage = ref(4)

if (headers.length === 0) {
  headers.push('id', 'name')
}

const pages = computed(() => Math.ceil(items.length / itemsPerPage.value))
const currentPage = ref(1)
const sortBy = ref(headers[0]!)
const sortOrder = ref(1)

const itemsSorted = computed(() => {
  if (!headers.includes(sortBy.value)) {
    return items
  }
  return items.slice().sort((a, b) => {
    const result = a[sortBy.value as keyof T] > b[sortBy.value as keyof T] ? 1 : -1
    return sortOrder.value * result
  })
})

const itemsDisplayed = computed(() => {
  const from = (currentPage.value - 1) * itemsPerPage.value
  const to = from + itemsPerPage.value
  return itemsSorted.value.slice(from, to)
})

function handlePageChange(selectedPage: number) {
  currentPage.value = selectedPage
}

function sortItemsBy(item: string) {
  const sorter = headers.find((h) => h === item) ?? 'id'
  sortOrder.value = sortBy.value === sorter ? -sortOrder.value : 1
  sortBy.value = sorter
}
</script>

<template>
  <div class="m-4">
    <div class="overflow-x-auto">
      <table class="m-4 w-full">
        <thead>
          <tr>
            <template v-for="header in headers">
              <th>
                <button
                  class="cursor-pointer p-2 hover:bg-gray-200"
                  :class="{ border: header === sortBy }"
                  @click="sortItemsBy(header)"
                  :aria-sort="
                    header === sortBy ? (sortOrder === 1 ? 'ascending' : 'descending') : 'none'
                  "
                >
                  <span class="mr-1">{{ header.toLocaleUpperCase('de') }}</span
                  ><span v-if="header !== sortBy">↕</span><span v-else-if="sortOrder === 1">↑</span
                  ><span v-else>↓</span>
                </button>
              </th>
            </template>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="item in itemsDisplayed">
            <TableRow :item="item" :hide-default="headers.length > 2">
              <template #extraCols="{ item }">
                <slot name="extraCols" :item="item"></slot>
              </template>
            </TableRow>
          </template>
        </tbody>
      </table>
    </div>
    <Pagination v-model:pages="pages" @page-change="handlePageChange" />
  </div>
</template>
