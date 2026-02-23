<template>
  <div class="flex flex-col items-start gap-1 flex-shrink-0" role="toolbar" aria-label="Block actions">
    <div class="flex items-center gap-0.5">
      <button type="button" class="btn btn-ghost btn-xs btn-square" title="Insert above" aria-label="Insert block above"
        @click="onInsertAbove">
        <IconInsertAbove :size="20" class="w-5 h-5" />
      </button>
      <button type="button" class="btn btn-ghost btn-xs btn-square" title="Insert below" aria-label="Insert block below"
        @click="onInsertBelow">
        <IconInsertBelow :size="20" class="w-5 h-5" />
      </button>
      <button v-if="isSection" type="button" class="btn btn-ghost btn-xs btn-square" title="Insert nested"
        aria-label="Insert block nested" @click="onInsertNest">
        <IconInsertNestBelow :size="20" class="w-5 h-5" />
      </button>
      <button type="button" class="btn btn-ghost btn-xs btn-square opacity-50 cursor-not-allowed"
        title="Move up (not implemented)" aria-label="Move up" disabled>
        <IconMoveUp :size="20" class="w-5 h-5" />
      </button>
      <button type="button" class="btn btn-ghost btn-xs btn-square opacity-50 cursor-not-allowed"
        title="Move down (not implemented)" aria-label="Move down" disabled>
        <IconMoveDown :size="20" class="w-5 h-5" />
      </button>
      <button type="button" class="btn btn-ghost btn-xs btn-square text-error hover:bg-error/10" title="Delete"
        aria-label="Delete block" @click="onDelete">
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

defineProps<{
  isSection: boolean
  isDirty?: boolean
}>()

const emit = defineEmits<{
  insertAbove: []
  insertBelow: []
  insertNest: []
  confirm: []
  cancel: []
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
function onDelete() {
  emit('delete')
}
</script>
