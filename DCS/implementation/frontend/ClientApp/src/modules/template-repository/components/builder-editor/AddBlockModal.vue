<template>
  <Teleport to="body">
    <div v-if="addBlockModalContext !== null" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      role="dialog" aria-modal="true" aria-labelledby="add-block-title" @click.self="handleCancel">
      <div class="bg-base-100 rounded-2xl shadow-xl w-full max-w-md mx-4 flex flex-col gap-4 p-6" @click.stop>
        <h2 id="add-block-title" class="text-lg font-bold">Add block</h2>
        <p class="text-sm text-base-content/70">Choose a block type:</p>
        <div class="flex flex-col gap-2">
          <BlockPaletteItem v-for="item in paletteBlockTypes" :key="item.blockType" :label="item.label"
            @select="handleAddBlock(item.blockType)" />
        </div>
        <div class="flex justify-end pt-2">
          <button type="button" class="btn btn-ghost btn-sm" @click="handleCancel">Cancel</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useTemplateDraftStore } from '@template-repository/store/templateDraftStore'
import { useTemplateEditorUiStore } from '@template-repository/store/templateEditorUiStore'
import { DocumentBlockType } from '@template-repository/models/contract-templace'
import BlockPaletteItem from '@template-repository/components/builder-editor/document-block/BlockPaletteItem.vue'

const draftStore = useTemplateDraftStore()
const uiStore = useTemplateEditorUiStore()
const { addBlockModalContext } = storeToRefs(uiStore)

/** TODO: add Clause when ready. */
const paletteBlockTypes = [
  { blockType: DocumentBlockType.Section, label: 'Section' },
  { blockType: DocumentBlockType.Text, label: 'Text' },
] as const

function handleCancel() {
  uiStore.closeAddBlockModal()
}

function handleAddBlock(blockType: typeof DocumentBlockType[keyof typeof DocumentBlockType]) {
  const ctx = addBlockModalContext.value
  if (ctx === null) return

  // After add, the new block's input will be auto-focused for editing content.
  draftStore.addBlock(ctx.parentBlockId, ctx.insertIndex, { blockType, text: '' })
  uiStore.closeAddBlockModal()
}
</script>
