<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  search: [value: string]
}>()

const search = ref('')

function searchList() {
  if (search.value.length < 1) return

  console.log(`Searching: ${search.value}`)
  emit('search', search.value)
}
</script>

<template>
  <div class="join m-2">
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
