import { useAuthStore } from '@/stores/auth'

export const AuthenticationService = {
  async login(username: string, password: string) {
    console.log('Signing in...')
    const success = Math.ceil(Math.random() * 3) > 1

    if (success) {
      const authStore = useAuthStore()
      authStore.setUser(username)
    }
    return success
  },

  async logout() {
    console.log('Signing out...')
    const authStore = useAuthStore()
    authStore.remove()
  },
}
