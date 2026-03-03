import authHttp from '@/api/auth-http'
import type { AuthCallbackResponse } from '@/models/responses/auth-callback-response'
import type { LoginResponse } from '@/models/responses/login-response'
import { useAuthStore } from '@/stores/auth-store'
import { useAuthTokenStore } from '@/stores/auth-token-store'

export const AuthenticationService = {
  async getLoginPath() {
    return await authHttp
      .get<LoginResponse>('/auth/login')
      .then((res) => res.data.auth_url)
      .catch((err) => {
        console.error('Login Error:', err)
        return ''
      })
  },

  async refresh() {
    return authHttp
      .post<AuthCallbackResponse>('/auth/refresh')
      .then((res) => {
        const authTokenStore = useAuthTokenStore()
        const resp = res.data
        authTokenStore.setTokens(resp.token_type, resp.access_token)
        const authStore = useAuthStore()
        authStore.setUser(resp.access_token)
        return res.data
      })
      .catch((err) => {
        if (err && err.status === 401) {
          const authStore = useAuthStore()
          authStore.remove()
          const authTokenStore = useAuthTokenStore()
          authTokenStore.remove()
        }
      })
  },

  async logout() {
    try {
      const response = await authHttp.get<{ logout_url: string }>('/auth/logout')
      const logoutUrl = response.data.logout_url
      
      // Clear local state first (but keep cookie for now)
      const authStore = useAuthStore()
      authStore.remove()
      const authTokenStore = useAuthTokenStore()
      authTokenStore.remove()
      
      // Redirect to Keycloak logout
      window.location.href = logoutUrl
    } catch (err) {
      console.error('Logout Error:', err)
      // Fallback: clear local state and redirect to home
      const authStore = useAuthStore()
      authStore.remove()
      const authTokenStore = useAuthTokenStore()
      authTokenStore.remove()
      window.location.href = '/'
    }
  },

  async logoutComplete() {
    try {
      await authHttp.get('/auth/logout-complete')
    } catch (err) {
      console.error('Logout complete error:', err)
    }
  },
}
