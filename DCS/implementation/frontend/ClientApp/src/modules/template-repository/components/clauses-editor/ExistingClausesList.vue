<template>
  <div class="space-y-2">
    <p v-if="!clauseBlocks.length" class="text-center py-6 text-xs text-base-content/40 italic">
      No clauses defined yet.
    </p>
    <div v-for="clause in clauseBlocks" :key="clause.blockId"
      class="flex items-start gap-3 p-3 rounded-lg border border-base-300 bg-base-200/30 group hover:shadow-sm transition-all">
      <div class="flex-1 min-w-0">
        <div class="font-semibold text-sm text-base-content">
          {{ clause.title ?? "" }}
        </div>
        <p class="text-xs text-base-content/70 mt-1 leading-relaxed whitespace-pre-wrap">
          <template v-for="(seg, i) in getSegments(clause)" :key="i">
            <template v-if="isText(seg)">{{ seg.value }}</template>
            <span v-else-if="isPlaceholder(seg)"
              class="clause-placeholder-slot inline-block border-b-2 border-dashed border-primary/50 bg-primary/5 text-primary px-1 rounded-sm font-medium cursor-default">
              {{ getPlaceholderLabel(seg) }}
            </span>
          </template>
        </p>
      </div>
      <button type="button"
        class="btn btn-ghost btn-xs text-error opacity-0 group-hover:opacity-100 transition-opacity shrink-0"
        aria-label="Delete clause" @click="$emit('delete', clause.blockId)"> ✕ </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ClauseBlock, SemanticCondition } from '@template-repository/models/contract-templace'
import { parseSegments, isText, isPlaceholder, type Segment } from '@template-repository/composables/useClauseTextChips'

const props = defineProps<{
  clauseBlocks: ClauseBlock[]
  semanticConditions: SemanticCondition[]
  getConditionName: (conditionId: string) => string
}>()

defineEmits<{
  delete: [blockId: string]
}>()

function getSegments(clause: ClauseBlock): Segment[] {
  return parseSegments(clause.text ?? '', props.semanticConditions)
}

function getParamType(conditionId: string, parameterName: string): string {
  const cond = props.semanticConditions.find((c) => c.conditionId === conditionId)
  const param = cond?.parameters.find((p) => p.parameterName === parameterName)
  return param?.type ?? 'string'
}

function getPlaceholderLabel(seg: Segment): string {
  if (!isPlaceholder(seg)) return ''
  const t = getParamType(seg.conditionId, seg.parameterName)
  return `${seg.parameterName} (${t})`
}
</script>
