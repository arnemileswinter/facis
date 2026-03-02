import type { AuthCallbackRequest } from '@/models/requests/auth-callback-request'
import type { AuthCallbackResponse } from '@/models/responses/auth-callback-response'
import type { LoginResponse } from '@/models/responses/login-response'
import { useAuthStore } from '@/stores/auth-store'
import axios from 'axios'

const http = axios.create({
  baseURL: import.meta.env.DCS_API_URL,
  headers: { 'Content-Type': 'application/json' },
})

export const AuthenticationService = {
  async getLoginPath() {
    return await http
      .get<LoginResponse>('/auth/login')
      .then((res) => res.data.auth_url)
      .catch((err) => {
        console.error('Login Error:', err)
        return ''
      })
  },

  async callback(request: AuthCallbackRequest) {
    return http
      .get<AuthCallbackResponse>('/auth/callback', { params: { ...request } })
      .then((res) => {
        const resp = res.data
        localStorage.setItem('access_token', resp.access_token)
        localStorage.setItem('token_type', resp.token_type)
        const authStore = useAuthStore()
        authStore.setUser(resp.access_token)
        return res.data
      })
      .catch((err) => {
        if (err && err.status === 401) {
          console.log(err)
          this.refresh()
        }
      })
  },

  async refresh() {
    return http
      .post<AuthCallbackResponse>('/auth/refresh')
      .then((res) => {
        return res.data
      })
      .catch((err) => {
        if (err && err.status === 401) {
          const authStore = useAuthStore()
          authStore.remove()
        }
      })
  },
}
