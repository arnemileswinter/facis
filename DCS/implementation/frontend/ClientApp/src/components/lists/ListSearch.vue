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
    <input
      type="text"
      v-model="search"
      @keyup.enter="searchList"
      placeholder="Search templates"
      aria-label="Search"
      class="input join-item"
    />
    <button @click="searchList" class="btn join-item">Search</button>
  </div>
</template>
