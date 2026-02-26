<template>
  <div class="space-y-6">
    <!-- Section 1: New clause -->
    <section class="rounded-lg border border-base-300 bg-base-100 p-4 shadow-sm">
      <h3 class="text-sm font-semibold text-base-content/80 mb-4">New clause</h3>
      <div class="space-y-4">
        <div>
          <label class="label-text text-xs text-base-content/60 block mb-1">Clause title
            <RequiredIndicator />
          </label>
          <input v-model="newClause.title" type="text" class="input input-bordered input-sm w-full" placeholder=""
            required />
        </div>
        <div>
          <label class="label-text text-xs text-base-content/60 block mb-1">Clause text
            <RequiredIndicator />
          </label>
          <ClauseTextEditor :model-value="newClause.text" :semantic-conditions="semanticConditions"
            @update:model-value="newClause.text = $event" />
        </div>
        <div class="flex justify-end">
          <button type="button" class="btn btn-secondary btn-sm" :disabled="!canAddClause" @click="addClause">
            Add clause
          </button>
        </div>
      </div>
    </section>

    <!-- Section 2: Existing clauses -->
    <section class="rounded-lg border border-base-300 bg-base-100 p-4 shadow-sm">
      <h3 class="text-sm font-semibold text-base-content/80 mb-4">Existing clauses</h3>
      <ExistingClausesList :clause-blocks="clauseBlocks" :semantic-conditions="semanticConditions"
        :get-condition-name="getConditionName" @delete="deleteClause" />
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useTemplateDraftStore } from '@template-repository/store/templateDraftStore'
import { isClauseBlock, type ClauseBlock } from '@template-repository/models/contract-templace'
import RequiredIndicator from '@core/components/RequiredIndicator.vue'
import ClauseTextEditor from '@template-repository/components/clauses-editor/ClauseTextEditor.vue'
import ExistingClausesList from '@template-repository/components/clauses-editor/ExistingClausesList.vue'

const store = useTemplateDraftStore()
const { documentBlocks, semanticConditions } = storeToRefs(store)

const newClause = ref({ title: '', text: '' })

/** Extract conditionIds from clause text placeholders {{conditionId.parameterName}}. */
function conditionIdsFromText(text: string): string[] {
  const set = new Set<string>()
  const re = /\{\{([^}]+)\}\}/g
  let m: RegExpExecArray | null
  while ((m = re.exec(text)) !== null) {
    const inner = m[1] ?? ''
    const dot = inner.indexOf('.')
    const conditionId = dot >= 0 ? inner.slice(0, dot) : inner
    if (conditionId) set.add(conditionId)
  }
  return [...set]
}

const clauseBlocks = computed((): ClauseBlock[] =>
  documentBlocks.value.filter((b): b is ClauseBlock => isClauseBlock(b))
)

const canAddClause = computed(() => !!newClause.value.title?.trim() && !!newClause.value.text?.trim())

function getConditionName(conditionId: string): string {
  const c = semanticConditions.value.find((x) => x.conditionId === conditionId)
  return c?.conditionName ?? conditionId
}

function addClause() {
  const text = newClause.value.text?.trim()
  if (!text) return
  store.addClause({
    title: newClause.value.title?.trim() || undefined,
    text,
    conditionIds: conditionIdsFromText(text),
  })
  newClause.value = { title: '', text: '' }
}

function deleteClause(blockId: string) {
  store.deleteClause(blockId)
}
</script>
