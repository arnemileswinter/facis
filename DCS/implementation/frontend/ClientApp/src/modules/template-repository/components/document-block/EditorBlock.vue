<template>
  <div :class="[
    'flex items-start gap-2 w-full rounded-lg border bg-base-100',
    'transition-[border-color,opacity] duration-200',
    borderClass
  ]" :data-block-id="item.blockId" @click="emit('select')" @focusin="emit('select')">
    <div class="flex-1 min-w-0 py-2 px-3">
      <!-- Section: title input -->
      <template v-if="block && isSectionBlock(block)">
        <label class="text-[10px] uppercase font-bold opacity-60">Section</label>
        <input v-model="localTitle" type="text" class="input input-sm input-ghost w-full font-semibold mt-0.5"
          placeholder="Section title" />
      </template>
      <!-- Text: textarea -->
      <template v-else-if="block && isTextBlock(block)">
        <label class="text-[10px] uppercase font-bold opacity-60">Text</label>
        <textarea v-model="localText"
          class="textarea textarea-ghost textarea-sm w-full mt-0.5 text-sm min-h-[2.5rem] resize-y"
          placeholder="Text content" rows="2" />
      </template>
      <div v-else class="rounded border border-base-300 bg-base-200/50 p-3 text-sm opacity-60">
        Block type "{{ block?.type }}" (TODO: Clause)
      </div>
    </div>
    <div class="pt-2 pr-2 pb-2">
      <BlockToolbar :item="item" :is-dirty="isDirty" @insert-above="emit('insertAbove')"
        @insert-below="emit('insertBelow')" @insert-nest="emit('insertNest')" @confirm="onConfirm"
        @cancel="revertToSaved" @move-up="emit('moveUp')" @move-down="emit('moveDown')"
        @move-outdent="emit('moveOutdent')" @move-indent="emit('moveIndent')" @delete="emit('delete')" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import type { EnrichedBlockItem } from '@template-repository/models/enriched-block-item'
import { isSectionBlock, isTextBlock } from '@template-repository/models/contract-templace'
import { useTemplateEditorUiStore } from '@template-repository/store/templateEditorUiStore'
import { useBlockMovementPreview } from '@template-repository/composables/useBlockMovementPreview'
import BlockToolbar from '@template-repository/components/toolbar/BlockToolbar.vue'

const props = defineProps<{
  item: EnrichedBlockItem
}>()

const emit = defineEmits<{
  select: []
  insertAbove: []
  insertBelow: []
  insertNest: []
  confirm: [payload: { title: string; text: string }]
  moveUp: []
  moveDown: []
  moveOutdent: []
  moveIndent: []
  delete: []
}>()

const uiStore = useTemplateEditorUiStore()
const { selectedBlockId } = storeToRefs(uiStore)
const { isSwapPreviewTarget } = useBlockMovementPreview()

const isSelected = computed(() => selectedBlockId.value === props.item.blockId)
const isSwapPreviewTargetForThis = computed(() => isSwapPreviewTarget(props.item.blockId))

const borderClass = computed(() => {
  if (isSwapPreviewTargetForThis.value) return 'border border-dashed border-primary'
  if (isSelected.value) return 'border border-primary'
  return 'border border-base-300'
})

const block = computed(() => props.item.block)

const savedTitle = computed(() => {
  const b = block.value
  if (b && isSectionBlock(b)) return b.title ?? b.text ?? ''
  return ''
})
const savedText = computed(() => {
  const b = block.value
  if (b && isTextBlock(b)) return b.text ?? ''
  if (b && isSectionBlock(b)) return b.text ?? ''
  return ''
})

const localTitle = ref('')
const localText = ref('')

watch(
  () => [props.item.blockId, savedTitle.value, savedText.value] as const,
  ([, title, text]) => {
    localTitle.value = title
    localText.value = text
  },
  { immediate: true }
)

const isDirty = computed(() => {
  const b = block.value
  if (b && isSectionBlock(b)) {
    return localTitle.value !== savedTitle.value || localText.value !== savedText.value
  }
  if (b && isTextBlock(b)) {
    return localText.value !== savedText.value
  }
  return false
})

function onConfirm() {
  emit('confirm', { title: localTitle.value, text: localText.value })
}

function revertToSaved() {
  localTitle.value = savedTitle.value
  localText.value = savedText.value
}
</script>
