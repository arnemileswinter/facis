import { defineStore } from 'pinia'
import { computed, ref, type Ref } from 'vue'
import { useAuthTokenStore } from './auth-token-store'
import { users } from '@/services/user-service'
import type { UserRole } from '@/types/user-role'

type User = string | null
interface UserI {
  id: string
  username: string
  name: string
  roles?: UserRole[]
}

export const useAuthStore = defineStore('auth', () => {
  const authTokenStore = useAuthTokenStore()
  const user: Ref<UserI | null> = ref(null)

  const isAuthenticated = computed(() => !!user.value && authTokenStore.isAuthSet)

  function setUser(newUser: User) {
    const userProfile = users.value.find((user) => user.id === newUser)
    if (!userProfile) return console.error('User Error: User not set')
    user.value = {
      id: userProfile.id,
      username: userProfile.username,
      name: userProfile.firstName + ' ' + userProfile.lastName,
      roles: userProfile.roleIds,
    }
  }

  function remove() {
    user.value = null
  }

  return { user, isAuthenticated, setUser, remove }
})
