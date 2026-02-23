<template>
  <div class="flex flex-col gap-2">
    <EditorBlock
      v-for="item in flatItemsWithBlock"
      :key="item.blockId"
      :block-id="item.blockId"
      :block="item.block"
      @insert-above="openAddBlockModal(item.resolvedParentBlockId, item.siblingIndex)"
      @insert-below="openAddBlockModal(item.resolvedParentBlockId, item.siblingIndex + 1)"
      @insert-nest="openAddBlockModal(item.blockId, 0)"
      @confirm="(payload) => confirmBlock(item.blockId, payload)"
      @delete="deleteBlock(item.blockId)"
    />
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

const rootBlockId = computed(() => documentOutline.value.find((b) => b.isRoot)?.blockId)

const flatItemsWithBlock = computed(() => {
  const list = flattened.value
  const blocks = documentBlocks.value
  const rootId = rootBlockId.value
  const blockById = new Map(blocks.map((b) => [b.blockId, b]))
  return list.map((item) => {
    const resolvedParentBlockId = item.parentBlockId === 'root' ? rootId ?? item.parentBlockId : item.parentBlockId
    return {
      blockId: item.blockId,
      block: blockById.get(item.blockId),
      siblingIndex: item.siblingIndex,
      resolvedParentBlockId: resolvedParentBlockId ?? item.parentBlockId,
    }
  })
})

function openAddBlockModal(parentBlockId: string, insertIndex: number) {
  uiStore.openAddBlockModal(parentBlockId, insertIndex)
}
function confirmBlock(blockId: string, payload: { title: string; text: string }) {
  draftStore.updateBlock(blockId, payload)
}
function deleteBlock(blockId: string) {
  draftStore.deleteBlock(blockId)
}
</script>
