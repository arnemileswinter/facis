import { defineStore } from 'pinia'
import { computed, ref, type Ref } from 'vue'

type User = string | null

export const useAuthStore = defineStore('auth', () => {
  const user: Ref<User> = ref(null)

  const isAuthenticated = computed(() => !!user.value)

  function setUser(newUser: User) {
    user.value = newUser
  }

  function remove() {
    user.value = null
  }

  return { user, isAuthenticated, setUser, remove }
})
