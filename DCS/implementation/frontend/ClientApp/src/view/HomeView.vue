<script setup lang="ts">
import packageJson from '../../package.json'
</script>
<template>
    <div>Hallo {{ packageJson.version }}</div>
</template>