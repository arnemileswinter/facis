<template>
  <div class="flex flex-col items-start gap-1 flex-shrink-0" role="toolbar" aria-label="Block actions">
    <div class="flex items-center gap-0.5">
      <button type="button" :class="btnIcon" title="Insert above" aria-label="Insert block above"
        @click="onInsertAbove">
        <IconInsertAbove :size="20" class="w-5 h-5" />
      </button>
      <button type="button" :class="btnIcon" title="Insert below" aria-label="Insert block below"
        @click="onInsertBelow">
        <IconInsertBelow :size="20" class="w-5 h-5" />
      </button>
      <button v-if="isSection" type="button" :class="btnIcon" title="Insert nested" aria-label="Insert block nested"
        @click="onInsertNest">
        <IconInsertNestBelow :size="20" class="w-5 h-5" />
      </button>
      <button type="button" :class="[btnIcon, !canMoveUp && 'opacity-50 cursor-not-allowed']" title="Move up"
        aria-label="Move up" :disabled="!canMoveUp" @click="onMoveUp">
        <IconMoveUp :size="20" class="w-5 h-5" />
      </button>
      <button type="button" :class="[btnIcon, !canMoveDown && 'opacity-50 cursor-not-allowed']" title="Move down"
        aria-label="Move down" :disabled="!canMoveDown" @click="onMoveDown">
        <IconMoveDown :size="20" class="w-5 h-5" />
      </button>
      <button type="button" :class="[btnIcon, 'text-error hover:bg-error/10']" title="Delete" aria-label="Delete block"
        @click="onDelete">
        <IconTrash :size="20" class="w-5 h-5" />
      </button>
    </div>
    <div v-if="isDirty" class="flex items-center gap-1 w-full">
      <button type="button" class="btn btn-ghost btn-xs flex-1" @click="onCancel">
        Cancel
      </button>
      <button type="button" class="btn btn-primary btn-xs flex-1" @click="onConfirm">
        Confirm
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import IconInsertAbove from '@template-repository/components/toolbar/icons/IconInsertAbove.vue'
import IconInsertBelow from '@template-repository/components/toolbar/icons/IconInsertBelow.vue'
import IconInsertNestBelow from '@template-repository/components/toolbar/icons/IconInsertNestBelow.vue'
import IconTrash from '@template-repository/components/toolbar/icons/IconTrash.vue'
import IconMoveUp from '@template-repository/components/toolbar/icons/IconMoveUp.vue'
import IconMoveDown from '@template-repository/components/toolbar/icons/IconMoveDown.vue'

const btnIcon = 'btn btn-ghost btn-xs btn-square'

defineProps<{
  isSection: boolean
  isDirty?: boolean
  canMoveUp?: boolean
  canMoveDown?: boolean
}>()

const emit = defineEmits<{
  insertAbove: []
  insertBelow: []
  insertNest: []
  confirm: []
  cancel: []
  moveUp: []
  moveDown: []
  delete: []
}>()

function onInsertAbove() {
  emit('insertAbove')
}
function onInsertBelow() {
  emit('insertBelow')
}
function onInsertNest() {
  emit('insertNest')
}
function onConfirm() {
  emit('confirm')
}
function onCancel() {
  emit('cancel')
}
function onMoveUp() {
  emit('moveUp')
}
function onMoveDown() {
  emit('moveDown')
}
function onDelete() {
  emit('delete')
}
</script>
