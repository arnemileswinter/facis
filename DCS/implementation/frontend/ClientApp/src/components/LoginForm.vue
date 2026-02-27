<script setup lang="ts">
import { ref } from 'vue'
import AuthButton from './AuthButton.vue'

const loginError = ref(false)

</script>

<template>
  <div class="bg-base-300 m-4 border-base-100 rounded-box w-[50%] border p-4 pt-0">
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-2xl">Login</legend>
      <div class="text-lg">Login via</div>
      <AuthButton />
    </fieldset>
  </div>

  <div v-if="loginError" class="toast toast-center">
    <div class="alert alert-error">Login failed!</div>
  </div>
</template>
