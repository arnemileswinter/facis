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
          <span v-if="outlineBlockIds.has(clause.blockId)" class="font-normal text-base-content/60 ml-1">
            (used in builder)
          </span>
        </div>
        <p class="text-xs text-base-content/70 mt-1 leading-relaxed whitespace-pre-wrap">
          <template v-for="(seg, i) in getSegments(clause)" :key="i">
            <template v-if="isText(seg)">{{ seg.value }}</template>
            <ClausePlaceholderSpan v-else-if="isPlaceholder(seg)" :label="getPlaceholderLabel(seg)" />
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
import { computed } from 'vue'
import type { ClauseBlock, SemanticCondition } from '@template-repository/models/contract-templace'
import { parseSegments, isText, isPlaceholder, type Segment } from '@template-repository/composables/useClauseTextChips'
import ClausePlaceholderSpan from '@template-repository/components/clauses-editor/ClausePlaceholderSpan.vue'

const props = defineProps<{
  clauseBlocks: ClauseBlock[]
  semanticConditions: SemanticCondition[]
  getConditionName: (conditionId: string) => string
  blockIdsInOutline: Set<string>
}>()

const outlineBlockIds = computed(() => props.blockIdsInOutline)

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
