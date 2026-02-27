<template>
  <span class="text-sm text-base-content leading-relaxed whitespace-nowrap">
    {{ text }}
  </span>
</template>

<script setup lang="ts">
const props = defineProps<{
  text: string
}>()
</script>
