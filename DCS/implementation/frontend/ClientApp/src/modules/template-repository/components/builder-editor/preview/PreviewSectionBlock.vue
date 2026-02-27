<template>
  <section class="mb-4">
    <h2 class="font-bold border-b border-base-300 pb-1 mb-1 text-base-content" :class="headingClass">
      {{ title }}
    </h2>
    <div v-if="hasChildren" class="pl-0 flex flex-wrap items-end gap-x-2 gap-y-1 section-children">
      <slot />
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
const props = defineProps<{
  title: string
  hasChildren?: boolean
  level?: number
}>()

const headingClass = computed(() => {
  const level = props.level ?? 1
  if (level <= 1) return 'text-lg'        // top-level section
  if (level === 2) return 'text-sm'       // second level, slightly smaller
  return 'text-xs'                        // deeper levels
})
</script>

<style scoped>
.section-children>section {
  width: 100%;
  flex-basis: 100%;
}
</style>
