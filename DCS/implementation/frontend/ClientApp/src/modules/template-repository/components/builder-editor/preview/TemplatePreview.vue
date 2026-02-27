<template>
  <div v-if="!hasBlockId">
    <template v-for="id in rootChildren" :key="id">
      <TemplatePreview :block-id="id" :section-level="1" />
    </template>
  </div>
  <template v-else>
    <PreviewSectionBlock v-if="block && isSection" :title="sectionTitle" :has-children="childrenIds.length > 0"
      :level="sectionLevel">
      <template v-for="childId in childrenIds" :key="childId">
        <TemplatePreview :block-id="childId" :section-level="sectionLevel + 1" />
      </template>
    </PreviewSectionBlock>
    <PreviewTextBlock v-else-if="block && isText" :text="block.text ?? ''" />
    <PreviewClauseBlock v-else-if="block && isClause" :text="block.text ?? ''"
      :semantic-conditions="semanticConditions" />
  </template>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useTemplateDraftStore } from '@template-repository/store/templateDraftStore'
import type { DocumentBlock, SectionBlock } from '@template-repository/models/contract-templace'
import {
  isSectionBlock,
  isTextBlock,
  isClauseBlock,
} from '@template-repository/models/contract-templace'
import PreviewSectionBlock from './PreviewSectionBlock.vue'
import PreviewTextBlock from './PreviewTextBlock.vue'
import PreviewClauseBlock from './PreviewClauseBlock.vue'

const props = defineProps<{
  /** If blockId is provided, the preview will render only that block and its children.
   *  If not provided, it will render all root-level blocks.
   */
  blockId?: string
  /** Section nesting level for headings (1 = top-level) */
  sectionLevel?: number
}>()

const draftStore = useTemplateDraftStore()
const { documentOutline, documentBlocks, semanticConditions } = storeToRefs(draftStore)

const hasBlockId = computed(() => props.blockId != null)

const rootChildren = computed(() => {
  const root = documentOutline.value.find((b) => b.isRoot)
  return root?.children ?? []
})

const block = computed<DocumentBlock | undefined>(() => {
  if (!props.blockId) return undefined
  return documentBlocks.value.find((b) => b.blockId === props.blockId)
})

const outlineNode = computed(() => {
  if (!props.blockId) return undefined
  return documentOutline.value.find((o) => o.blockId === props.blockId)
})

const childrenIds = computed(() => outlineNode.value?.children ?? [])

const isSection = computed(() => !!block.value && isSectionBlock(block.value))
const isText = computed(() => !!block.value && isTextBlock(block.value))
const isClause = computed(() => !!block.value && isClauseBlock(block.value))

const sectionTitle = computed(() => {
  const b = block.value as SectionBlock | undefined
  return b?.title ?? b?.text ?? ''
})

const sectionLevel = computed(() => props.sectionLevel ?? 1)
</script>
