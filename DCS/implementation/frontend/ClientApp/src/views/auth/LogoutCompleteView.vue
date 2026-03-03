<script setup lang="ts">
import { AuthenticationService } from '@/services/authentication-service'
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

onMounted(async () => {
  // Finalize logout by clearing the cookie on the backend
  await AuthenticationService.logoutComplete()
  // Redirect to home/login page
  router.replace({ name: 'home' })
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-base-200">
    <span class="loading loading-spinner loading-lg" />
  </div>
</template>
