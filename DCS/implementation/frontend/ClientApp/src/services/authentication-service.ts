import { useAuthStore } from '@/stores/auth-store'
import axios from 'axios'

const http = axios.create({
  baseURL: import.meta.env.DCS_API_URL,
  headers: { 'Content-Type': 'application/json' },
})

export const AuthenticationService = {
  async getLoginPath() {
     const login = 'https://keycloak.xfsc.local/realms/dcs/protocol/openid-connect/auth?client_id=digital-contracting-service&redirect_uri=http%3A%2F%2Flocalhost%3A8991%2Fauth%2Fcallback&response_type=code&scope=openid'
    const result = await http
      .get<{ auth_url: string }>('/auth/login')
      .then((res) => res.data.auth_url)
      .catch((err) => {
        console.error('Login Error:', err)
        return ''
      })
    return result ? result : login
  },

  async callback(code: string) {
    return http.get('/auth/callback', {params: {code: code}}).then(res => res.data)
  },

  async refresh() {
    return http.post('/auth/refresh').then(res => res.data)
  },

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
