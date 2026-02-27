<script setup lang="ts">
import type { AuthCallbackRequest } from '@/models/requests/auth-callback-request'
import { AuthenticationService } from '@/services/authentication-service'
import { onMounted, ref, type Ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const query = route.query

const callbackQuery: Ref<AuthCallbackRequest | null> = ref(null)
if (
  query.iss &&
  !Array.isArray(query.iss) &&
  query.session_state &&
  !Array.isArray(query.session_state) &&
  query.code &&
  !Array.isArray(query.code)
) {
  callbackQuery.value = {
    iss: query.iss,
    code: query.code,
    session_state: query.session_state,
  }
}

onMounted(async () => {
  if (callbackQuery.value) {
    await AuthenticationService.callback(callbackQuery.value)
  }
  router.replace({ name: 'templates.list' })
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-base-200">
    <span class="loading loading-spinner loading-lg" />
  </div>
</template>
