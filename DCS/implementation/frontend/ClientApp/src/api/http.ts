import { AuthenticationService } from '@/services/authentication-service'
import { useAuthTokenStore } from '@/stores/auth-token-store'
import axios from 'axios'
import { storeToRefs } from 'pinia'

const API_BASE_URL = import.meta.env.DCS_API_URL

const http = axios.create({
  baseURL: API_BASE_URL,
  headers: { 'Content-Type': 'application/json' },
})

http.interceptors.request.use(
  (config) => {
    const tokenStore = useAuthTokenStore()
    const { isAuthSet, getAuthenticationHeader } = storeToRefs(tokenStore)
    if (isAuthSet.value) {
      config.headers.Authorization = getAuthenticationHeader.value
    } else {
      delete config.headers.Authorization
    }
    return config
  },
  (err) => Promise.reject(err),
)

http.interceptors.response.use(
  (resp) => resp,
  async (err) => {
    if (err.status === 401) {
      const res = await AuthenticationService.refresh()
      if (res) {
        return http(err.config)
      }
    }
    return Promise.reject(err)
  },
)

export default http
