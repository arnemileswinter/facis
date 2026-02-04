<script setup lang="ts" generic="T extends TableItem">
import { computed, ref, watchEffect } from 'vue'
import type { TableItem } from '../models/table-item'
import Pagination from './Pagination.vue'
import TableRow from './TableRow.vue'

const { items } = defineProps<{
  items: T[]
}>()

const itemsPerPage = ref(2)

const headers = ['id', 'name']
const pages = ref(0)
const page = ref(1)
const sortBy = ref('id')
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
  const from = (page.value - 1) * itemsPerPage.value
  const to = from + itemsPerPage.value
  return itemsSorted.value.slice(from, to)
})

function changePage(selectedPage: number) {
  page.value = selectedPage
}

function sortItemsBy(item: string) {
  const sorter = headers.find((h) => h === item) ?? 'id'
  sortOrder.value = sortBy.value === sorter ? -sortOrder.value : 1
  sortBy.value = sorter
}

watchEffect(() => {
  pages.value = Math.ceil(items.length / itemsPerPage.value)
})
</script>

<template>
  <div class="m-4">
    <div class="overflow-x-auto">
      <table class="m-4">
        <thead>
          <tr>
            <template v-for="header in headers">
              <th>
                <button
                  class="cursor-pointer p-2"
                  :class="{ border: header === sortBy }"
                  @click="sortItemsBy(header)"
                  :aria-sort="
                    header === sortBy ? (sortOrder === 1 ? 'ascending' : 'descending') : 'none'
                  "
                >
                  <span>{{ header.toLocaleUpperCase('de') }}</span
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
            <TableRow :item="item">
              <template #extraRows="{ item }"></template>
            </TableRow>
          </template>
        </tbody>
      </table>
    </div>
    <Pagination v-model:pages="pages" @change="changePage" />
  </div>
</template>
