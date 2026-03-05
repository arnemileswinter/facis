<script setup lang="ts">
import type { PartialContractTemplate } from '@/models/contract-template'
import { ContractTemplateService } from '@/services/contract-template-service'
import { computed, ref, type Ref } from 'vue'

const props = defineProps<{
  items: PartialContractTemplate[]
}>()

const emit = defineEmits<{
  searchResult: [value: PartialContractTemplate[]]
}>()

const search = ref('')

const filterLabels = {
  did: 'DID',
  document_number: 'Document number',
  version: 'Version',
  template_type: 'Template type',
  state: 'State',
  name: 'Name',
  description: 'Description',
  filter: 'Filter',
} as const
type FilterLabels = typeof filterLabels
type FilterLabelKey = keyof FilterLabels
type FilterLabelValue = FilterLabels[FilterLabelKey]

const selectedFilter = ref<FilterLabelValue>('Name')

const searchResults: Ref<Set<string>> = ref(new Set())

const searchKey = computed(() => {
  return (Object.keys(filterLabels) as FilterLabelKey[]).find((key) => filterLabels[key] === selectedFilter.value)
})

const searchedItems = computed(() => {
  if (search.value.length < 1) return props.items
  return props.items.filter((item) => searchResults.value.has(`${item.did}|${item.document_number}|${item.version}`))
})

async function searchList() {
  if (search.value.length < 1 || !searchKey.value) {
    emit('searchResult', props.items)
    return
  }

  const request = { [searchKey.value]: search.value }
  const searchResult = await ContractTemplateService.search(request)
  searchResults.value = new Set(searchResult.map((item) => `${item.did}|${item.document_number}|${item.version}`))
  emit('searchResult', searchedItems.value)
}
</script>

<template>
  <div class="join m-2">
    <select v-model="selectedFilter" class="select select-neutral join-item w-max" aria-label="Search type">
      <option disabled>Select search filter</option>
      <option v-for="filter in filterLabels" :key="filter" :value="filter">
        {{ filter }}
      </option>
    </select>
    <label class="input input-neutral join-item grow">
      <input
        type="text"
        v-model="search"
        @keyup.enter="searchList"
        placeholder="Search templates"
        aria-label="Search templates"
      />
    </label>
    <button @click="searchList" class="btn btn-neutral join-item">Search</button>
  </div>
</template>
