<script setup lang="ts">
import { AuthenticationService } from '@/services/authentication-service'
import { useAuthStore } from '@/stores/auth-store'
import { storeToRefs } from 'pinia'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const authStore = useAuthStore()
const { isAuthenticated } = storeToRefs(authStore)

const loginPath = ref('')

onMounted(async () => {
  loginPath.value = await AuthenticationService.getLoginPath()
})

function logout() {
  AuthenticationService.logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <div>
    <button v-if="!isAuthenticated" class="btn btn-block btn-accent flex-1 text-center">
      <a :href="loginPath">Single Sign On</a>
    </button>
    <button v-else class="btn btn-block btn-accent flex-1 text-center" @click="logout">Logout</button>
  </div>
</template>
