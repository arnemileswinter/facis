<script setup lang="ts">
import { AuthenticationService } from '@/services/authentication-service'
import { ref } from 'vue'
import AuthButton from './AuthButton.vue'

const loginError = ref(false)

async function login() {
  const result = await AuthenticationService.login('test', 'test')
  if (!result) {
    loginError.value = true
    setTimeout(() => (loginError.value = false), 2000)
  }
}
</script>

<template>
  <form class="bg-base-300 m-4 border-base-100 rounded-box w-[50%] border p-4 pt-0" @submit.prevent="login">
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-2xl">Login</legend>
      <div class="text-lg">Login via</div>
      <AuthButton />
    </fieldset>
  </form>

  <div v-if="loginError" class="toast toast-center">
    <div class="alert alert-error">Login failed!</div>
  </div>
</template>
