<script setup lang="ts">
import { authenticationService } from '@/services/authentication-service'
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

onMounted(async () => {
  // Exchange refresh token cookie for access token
  const result = await authenticationService.refresh()
  // Redirect to templates list on success
  if (result) {
    router.replace({ name: 'templates.list' })
  } else {
    router.replace({ name: 'home' })
  }
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-base-200">
    <span class="loading loading-spinner loading-lg" />
  </div>
</template>
