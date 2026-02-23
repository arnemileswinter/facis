<template>
  <div class="flex flex-col gap-2">
    <div v-for="item in flatItemsWithBlock" :key="item.blockId" class="flex items-stretch min-w-0">
      <!-- Indent area: width by depth, left border for children to show hierarchy -->
      <div
        :class="['flex-shrink-0', 'transition-[width]', 'duration-300', 'ease-out', item.depthLevel > 0 && 'border-l-2 border-base-300']"
        :style="{ width: indentWidth(item.depthLevel) }" aria-hidden />
      <EditorBlock :block-id="item.blockId" :block="item.block" :can-move-up="item.siblingIndex > 0"
        :can-move-down="item.siblingIndex < item.siblingCount - 1" :can-outdent="item.canOutdent"
        :can-indent="item.canIndent" @insert-above="openAddBlockModal(item.parentBlockId, item.siblingIndex)"
        @insert-below="openAddBlockModal(item.parentBlockId, item.siblingIndex + 1)"
        @insert-nest="openAddBlockModal(item.blockId, 0)" @confirm="(payload) => confirmBlock(item.blockId, payload)"
        @move-up="moveBlockUp(item.blockId, item.parentBlockId, item.siblingIndex)"
        @move-down="moveBlockDown(item.blockId, item.parentBlockId, item.siblingIndex)"
        @move-outdent="moveBlockOutdent(item.blockId, item.outdentGrandparentBlockId, item.outdentInsertIndex)"
        @move-indent="moveBlockIndent(item.blockId, item.indentParentBlockId, item.indentInsertIndex)"
        @delete="deleteBlock(item.blockId)" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useTemplateDraftStore } from '@template-repository/store/templateDraftStore'
import { useTemplateEditorUiStore } from '@template-repository/store/templateEditorUiStore'
import {
  useFlattenedOutline,
  type FlattenedOutlineItem,
} from '@template-repository/composables/useFlattenedOutline'
import type { DocumentBlock, DocumentOutline, DocumentOutlineBlock } from '@template-repository/models/contract-templace'
import { isSectionBlock } from '@template-repository/models/contract-templace'
import EditorBlock from '@template-repository/components/document-block/EditorBlock.vue'

const draftStore = useTemplateDraftStore()
const uiStore = useTemplateEditorUiStore()
const { documentOutline, documentBlocks } = storeToRefs(draftStore)

const flattened = useFlattenedOutline(documentOutline)

const flatItemsWithBlock = computed(() => {
  const outline = documentOutline.value
  const root = outline.find((b) => b.isRoot)
  const blockById = new Map(documentBlocks.value.map((b) => [b.blockId, b]))
  return flattened.value.map((item) => enrichFlatItem(item, outline, blockById, root))
})

function openAddBlockModal(parentBlockId: string, insertIndex: number) {
  uiStore.openAddBlockModal(parentBlockId, insertIndex)
}
function confirmBlock(blockId: string, payload: { title: string; text: string }) {
  draftStore.updateBlock(blockId, payload)
}
function moveBlockUp(blockId: string, parentBlockId: string, siblingIndex: number) {
  draftStore.moveBlock(blockId, parentBlockId, siblingIndex - 1)
}
function moveBlockDown(blockId: string, parentBlockId: string, siblingIndex: number) {
  draftStore.moveBlock(blockId, parentBlockId, siblingIndex + 1)
}
function moveBlockOutdent(blockId: string, grandparentBlockId: string, insertIndex: number) {
  if (!grandparentBlockId) return
  draftStore.moveBlock(blockId, grandparentBlockId, insertIndex)
}
function moveBlockIndent(blockId: string, parentBlockId: string, insertIndex: number) {
  if (!parentBlockId) return
  draftStore.moveBlock(blockId, parentBlockId, insertIndex)
}
function deleteBlock(blockId: string) {
  draftStore.deleteBlock(blockId)
}

/**
 * Enriches a flattened outline item with block data and outdent/indent params for the toolbar.
 */
function enrichFlatItem(
  item: FlattenedOutlineItem,
  outline: DocumentOutline,
  blockById: Map<string, DocumentBlock>,
  root: DocumentOutlineBlock | undefined
) {
  const parentNode = outline.find((b) => b.blockId === item.parentBlockId)
  const siblingCount = parentNode?.children?.length ?? 0
  const isDirectChildOfRoot = !!root && item.parentBlockId === root.blockId
  const isLastChild = item.siblingIndex === siblingCount - 1
  const grandparentNode = outline.find((b) => b.children?.includes(item.parentBlockId))
  const parentIndexInGrandparent = grandparentNode?.children?.indexOf(item.parentBlockId) ?? -1
  const canOutdent = !isDirectChildOfRoot && isLastChild
  const outdentGrandparentBlockId = grandparentNode?.blockId ?? ''
  const outdentInsertIndex = parentIndexInGrandparent + 1

  const prevSiblingBlockId = parentNode?.children?.[item.siblingIndex - 1]
  const prevSiblingBlock = prevSiblingBlockId ? blockById.get(prevSiblingBlockId) : undefined
  const prevSiblingIsSection = !!prevSiblingBlock && isSectionBlock(prevSiblingBlock)
  const canIndent = item.siblingIndex > 0 && prevSiblingIsSection
  const prevSiblingOutlineNode = prevSiblingBlockId ? outline.find((b) => b.blockId === prevSiblingBlockId) : undefined
  const indentInsertIndex = prevSiblingOutlineNode?.children?.length ?? 0

  return {
    blockId: item.blockId,
    block: blockById.get(item.blockId),
    siblingIndex: item.siblingIndex,
    siblingCount,
    parentBlockId: item.parentBlockId,
    depthLevel: item.depthLevel,
    canOutdent,
    canIndent,
    outdentGrandparentBlockId,
    outdentInsertIndex,
    indentParentBlockId: prevSiblingBlockId ?? '',
    indentInsertIndex,
  }
}

/** Indent width per nesting level (px) */
const INDENT_PER_LEVEL = 16
function indentWidth(depth: number): string {
  return `${depth * INDENT_PER_LEVEL}px`
}
</script>
