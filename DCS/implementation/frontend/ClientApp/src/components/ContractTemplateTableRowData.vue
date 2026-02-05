<script setup lang="ts">
import type { ContractTemplate } from '../models/contract-template';

defineProps<{
    template: ContractTemplate
}>()
</script>

<template>
    <td class="pl-8">{{ template.did }}</td>
    <td class="pl-8">{{ template.documentNumber }}</td>
    <td class="pl-8">{{ template.name }}</td>
    <td class="pl-8">{{ template.state }}</td>
    <td class="pl-8">{{ template.version }}</td>
    <td class="pl-8">{{ template.createdBy }}</td>
    <td class="pl-8">{{ template.createdAt }}</td>
    <td class="pl-8">{{ template.description }}</td>
    <td class="pl-8">{{ template.meta_data }}</td>
</template>
