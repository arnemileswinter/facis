<script setup lang="ts">
import { AuthenticationService } from '@/services/authentication-service'
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

interface FormData {
  username: string
  password: string
}

const router = useRouter()
const route = useRoute()

const formData = reactive<FormData>({ username: '', password: '' })
const loginError = ref(false)

const rules = {
  username: { pattern: '.{3,}', info: 'Must be more than 3 characters' },
  password: { pattern: '.{4,}', info: 'Must be more than 4 characters' },
}

async function login() {
  const result = await AuthenticationService.login(formData.username, formData.password)
  if (!result) {
    loginError.value = true
    setTimeout(() => (loginError.value = false), 2000)
  } else {
    route.meta.hideInSidebar = true
    router.push({ name: 'templates.list' })
  }
}
</script>

<template>
  <form class="bg-base-300 m-4 border-base-100 rounded-box w-[50%] border p-4 pt-0" @submit.prevent="login">
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-2xl">Login</legend>

      <fieldset class="fieldset">
        <label class="label">User name</label>
        <input
          type="text"
          class="input validator w-full"
          required
          placeholder="User name"
          :pattern="rules.username.pattern"
          v-model="formData.username"
        />
        <p class="validator-hint">{{ rules.username.info }}</p>
      </fieldset>

      <fieldset class="fieldset">
        <label class="label">Password</label>
        <input
          type="password"
          class="input validator w-full"
          required
          placeholder="Password"
          :pattern="rules.password.pattern"
          v-model="formData.password"
        />
        <p class="validator-hint">{{ rules.password.info }}</p>
      </fieldset>

      <button type="submit" class="btn btn-block btn-accent">Login</button>
    </fieldset>
  </form>

  <div v-if="loginError" class="toast toast-center">
    <div class="alert alert-error">Login failed!</div>
  </div>
</template>
