<script setup lang="ts">
import { ref, watch } from 'vue'

const emit = defineEmits<{
  pageChange: [value: number]
}>()

const pages = defineModel<number>('pages', {
  required: true,
  set: (value) => (value >= 0 ? Math.floor(value) : 0),
})

const currentPage = ref(1)

watch(currentPage, (newPage, oldPage) => {
  if (newPage !== oldPage) {
    emit('pageChange', currentPage.value)
  }
})
</script>

<template>
  <div v-if="pages > 0" class="flex justify-center">
    <template v-for="page in pages">
      <button
        type="button"
        class="border p-2 my-0 cursor-pointer hover:bg-blue-200"
        :class="{ 'bg-amber-400': page == currentPage }"
        @click="currentPage = page"
      >
        {{ page }}
      </button>
    </template>
  </div>
</template>
