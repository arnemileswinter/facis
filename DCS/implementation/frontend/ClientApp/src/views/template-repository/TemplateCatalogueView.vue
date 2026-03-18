<template>
  <div>
    <div class="flex justify-between mb-8">
      <h2 class="text-2xl/7 font-bold sm:truncate sm:text-3xl sm:tracking-tight">
        {{ $route.meta.name ?? 'Template Catalogue' }}
      </h2>
    </div>

    <div>
      <div v-if="loading">Lade Templates...</div>
      <div v-else-if="error">{{ error }}</div>
      <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
        <!-- Details -->
        <div class="lg:col-span-2 px-2 sm:px-4 pb-6">
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0">
              <h3 class="text-lg font-semibold text-base-content truncate">
                {{ catalogue?.name || 'Template details' }}
              </h3>
            </div>
            <div v-if="catalogue?.templateType" class="badge badge-md badge-accent shrink-0">
              {{ String(catalogue.templateType) }}
            </div>
          </div>

          <template v-if="catalogue?.sdMeta" class="mt-4">
            <div class="mt-4">
              <div class="flex flex-col gap-3">
                <div class="form-control">
                  <label class="label">
                    <span class="label-text text-xs font-semibold tracking-wide">{{ toReadableName('sdHash') }}</span>
                  </label>
                  <input class="input input-bordered input-sm w-full font-mono text-xs" :value="catalogue.sdMeta.sdHash"
                    readonly />
                </div>
                <div class="form-control">
                  <label class="label">
                    <span class="label-text text-xs font-semibold tracking-wide">{{ toReadableName('issuer') }}</span>
                  </label>
                  <input class="input input-bordered input-sm w-full font-mono text-xs" :value="catalogue.sdMeta.issuer"
                    readonly />
                </div>
                <div class="form-control">
                  <label class="label">
                    <span class="label-text text-xs font-semibold tracking-wide">{{
                      toReadableName('uploadDatetime') }}
                    </span>
                  </label>
                  <input class="input input-bordered input-sm w-full font-mono text-xs"
                    :value="toReadableValue('uploadDatetime', catalogue.sdMeta.uploadDatetime)" readonly />
                </div>
              </div>
            </div>
            <div class="divider my-4"></div>
          </template>

          <div v-if="catalogue" class="flex flex-col gap-3">
            <template v-for="(value, key) in displayFields" :key="key">
              <div v-if="value" class="form-control">
                <label class="label">
                  <span class="label-text text-xs font-semibold tracking-wide">
                    {{ toReadableName(key) }}
                  </span>
                </label>
                <input class="input input-bordered input-sm w-full font-mono text-xs"
                  :value="toReadableValue(key, value)" readonly />
              </div>
            </template>
          </div>
          <p v-else class="text-sm text-base-content/70">No data loaded.</p>
        </div>

        <!-- Actions -->
        <div class="lg:pr-6">
          <div class="card bg-base-100 shadow-sm border border-base-300">
            <div class="card-body">
              <p class="text-sm text-base-content/70">
                This template is free to use.
              </p>
              <div class="pt-2 flex flex-col gap-2">
                <button class="btn btn-sm btn-primary rounded-box" disabled>
                  Get
                </button>
                <button class="btn btn-sm btn-ghost rounded-box" @click="router.back()">
                  Back
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>
  </div>

</template>
<script setup lang="ts">
import { useSelfDescriptionConverter } from '@/modules/template-catalogue/composables/useSelfDescriptionConverter';
import type { TemplateCatalogue } from '@/modules/template-catalogue/models/template-catalogue';
import { computed, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const router = useRouter()
const route = useRoute()
const { loading, error, loadTemplateCatalogue } = useSelfDescriptionConverter()
const catalogue = ref<TemplateCatalogue | null>(null)

const did = computed(() => `${route.params.did ?? ""}`)

const displayFields = computed(() => {
  if (!catalogue.value) return {}
  const { sdMeta, ...rest } = catalogue.value
  return rest as Record<string, unknown>
})

async function load() {
  catalogue.value = await loadTemplateCatalogue(did.value)
}
load()

function toReadableName(key: string): string {
  if (!key) return "";
  if (key.toLowerCase() === "did") return "DID"
  if (key.toLowerCase() === "sdhash") return "SD Hash"
  if (key.toLowerCase() === "uploaddatetime") return "Upload datetime"
  if (key.toLowerCase() === "statusdatetime") return "Status datetime"
  const words = key
    // split camelCase boundaries
    .replace(/([a-z0-9])([A-Z])/g, "$1 $2")
    // split cases like ABCWord → ABC Word
    .replace(/([A-Z])([A-Z][a-z])/g, "$1 $2")
    .split(" ");

  return words.map(word =>
    word.length ? word.charAt(0).toUpperCase() + word.slice(1).toLowerCase() : ""
  ).join(" ");
}

function toReadableValue(key: string, value: unknown): string {
  if (!value) return ""
  const lower = key.toLowerCase()
  if (["createdat", "updatedat", "uploaddatetime", "statusdatetime"].includes(lower)) {
    const d = new Date(String(value))
    return Number.isNaN(d.getTime()) ? String(value) : d.toLocaleDateString()
  }
  return String(value)
}

</script>