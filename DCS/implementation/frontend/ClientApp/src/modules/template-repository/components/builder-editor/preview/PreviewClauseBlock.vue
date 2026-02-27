<template>
  <span class="inline-flex items-baseline text-sm text-base-content leading-relaxed">
    <template v-for="(seg, index) in segments" :key="index">
      <span v-if="seg.type === 'text'"> {{ seg.value }} </span>
      <PreviewParamInput v-else :type="seg.paramType" :label="seg.label" />
    </template>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { SemanticCondition, SemanticParameterType } from '@template-repository/models/contract-templace'
import { parseSegments, isText, isPlaceholder, type Segment } from '@template-repository/composables/useClauseTextChips'
import PreviewParamInput from './PreviewParamInput.vue'

const props = defineProps<{
  text: string
  semanticConditions: SemanticCondition[]
}>()

type PreviewSegment =
  | { type: 'text'; value: string }
  | { type: 'param'; paramType: SemanticParameterType; label: string }

const segments = computed<PreviewSegment[]>(() => {
  const normalizedText = (props.text ?? '').replace(/^[\s\u00A0]+/, '')
  const baseSegments: Segment[] = parseSegments(normalizedText, props.semanticConditions)
  const result: PreviewSegment[] = []
  for (const seg of baseSegments) {
    if (isText(seg)) {
      result.push({ type: 'text', value: seg.value })
    } else if (isPlaceholder(seg)) {
      const cond = props.semanticConditions.find((c) => c.conditionId === seg.conditionId)
      const param = cond?.parameters.find((p) => p.parameterName === seg.parameterName)
      const paramType: SemanticParameterType = param?.type ?? 'string'
      result.push({
        type: 'param',
        paramType,
        label: seg.parameterName,
      })
    }
  }
  return result
})
</script>
