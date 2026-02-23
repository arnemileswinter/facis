<template>
  <div class="flex flex-col gap-2">
    <EditorBlock v-for="item in flatItemsWithBlock" :key="item.blockId" :block-id="item.blockId" :block="item.block"
      :can-move-up="item.siblingIndex > 0" :can-move-down="item.siblingIndex < item.siblingCount - 1"
      @insert-above="openAddBlockModal(item.parentBlockId, item.siblingIndex)"
      @insert-below="openAddBlockModal(item.parentBlockId, item.siblingIndex + 1)"
      @insert-nest="openAddBlockModal(item.blockId, 0)" @confirm="(payload) => confirmBlock(item.blockId, payload)"
      @move-up="moveBlockUp(item.blockId, item.parentBlockId, item.siblingIndex)"
      @move-down="moveBlockDown(item.blockId, item.parentBlockId, item.siblingIndex)"
      @delete="deleteBlock(item.blockId)" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useTemplateDraftStore } from '@template-repository/store/templateDraftStore'
import { useTemplateEditorUiStore } from '@template-repository/store/templateEditorUiStore'
import { useFlattenedOutline } from '@template-repository/composables/useFlattenedOutline'
import EditorBlock from '@template-repository/components/document-block/EditorBlock.vue'

const draftStore = useTemplateDraftStore()
const uiStore = useTemplateEditorUiStore()
const { documentOutline, documentBlocks } = storeToRefs(draftStore)

const flattened = useFlattenedOutline(documentOutline)

const flatItemsWithBlock = computed(() => {
  const list = flattened.value
  const blocks = documentBlocks.value
  const outline = documentOutline.value
  const blockById = new Map(blocks.map((b) => [b.blockId, b]))
  return list.map((item) => {
    const parentNode = outline.find((b) => b.blockId === item.parentBlockId)
    const siblingCount = parentNode?.children?.length ?? 0
    return {
      blockId: item.blockId,
      block: blockById.get(item.blockId),
      siblingIndex: item.siblingIndex,
      siblingCount: siblingCount,
      parentBlockId: item.parentBlockId,
    }
  })
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
function deleteBlock(blockId: string) {
  draftStore.deleteBlock(blockId)
}
</script>
