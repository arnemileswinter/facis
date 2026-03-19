<template>
  <div>
    <div class="flex justify-between mb-8">
      <h2 class="text-2xl/7 font-bold sm:truncate sm:text-3xl sm:tracking-tight">
        {{ $route.meta.name ?? 'Template Catalogues' }}
      </h2>
    </div>

    <div>
      <div v-if="loading">Lade Templates...</div>
      <div v-else-if="error">{{ error }}</div>
      <div v-else>
        <ul class="list">
          <li v-if="!items.length" class="text-center text-base-content/60 py-6">
            No data found.
          </li>
          <!-- card -->
          <li v-for="item in items" :key="item.did" class="list-row min-w-0 w-full">
            <div class="list-col-grow card bg-base-200 card-border hover:bg-base-300 min-w-0 w-full">
              <div class="card-body min-w-0">
                <h2 class="card-title justify-between">
                  <div class="flex gap-8 h-full">
                    <div>Name: {{ item.name }}</div>
                    <div v-if="item.templateType" class="badge badge-md badge-accent h-full">{{
                      toProperCase(item.templateType) }}</div>
                  </div>
                </h2>
                <div class="flex justify-between flex-col gap-2">
                  <div v-if="item.documentNumber">Document number: {{ item.documentNumber }}</div>
                  <div v-if="item.version">Version: {{ item.version }}</div>
                </div>
                <div class="flex justify-between min-w-0">
                  <div v-if="item.sdMeta?.uploadDatetime">
                    Upload Date: {{ new Date(item.sdMeta.uploadDatetime).toLocaleDateString() }}
                  </div>
                  <div class="card-actions justify-end flex-none">
                    <RouterLink :to="{ name: 'template.catalogues.view', params: { did: item.did } }"
                      class="btn btn-sm btn-primary rounded-box gap-2">
                      View
                    </RouterLink>
                  </div>
                </div>
              </div>
            </div>
          </li>
        </ul>

        <div v-if="totalPages > 1" class="mt-4 flex items-center justify-end gap-3">
          <span class="text-xs text-base-content/70">
            Page {{ page + 1 }} of {{ totalPages }}
          </span>
          <div class="join">
            <button class="btn btn-xs join-item" :disabled="!canPrev" @click="goPrev"> Prev </button>
            <button class="btn btn-xs join-item" :disabled="!canNext" @click="goNext"> Next </button>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { useSelfDescriptionConverter } from '@/modules/template-catalogue/composables/useSelfDescriptionConverter';
import type { TemplateCatalogue } from '@/modules/template-catalogue/models/template-catalogue';
import { toProperCase } from '@/utils/string';
import { computed, ref } from 'vue';
const pageSize = 20
const page = ref(0)
const items = ref<TemplateCatalogue[]>([])

const { loading, error, totalCount, loadTemplateCatalogues } = useSelfDescriptionConverter()
async function load() {
  items.value = await loadTemplateCatalogues(page.value, pageSize)
}
load()

const totalPages = computed(() => {
  if (loading.value) return 0
  return totalCount.value > 0 ? Math.ceil(totalCount.value / pageSize) : 1
})

const canPrev = computed(() => page.value > 0)
const canNext = computed(() => (page.value + 1) * pageSize < totalCount.value)

function goPrev() {
  if (!canPrev.value) return
  page.value -= 1
  load()
}

function goNext() {
  if (!canNext.value) return
  page.value += 1
  load()
}
</script>