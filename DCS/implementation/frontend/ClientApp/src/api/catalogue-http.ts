import axios from 'axios'
import { getConfig } from '@/config'
import { useAuthTokenStore } from '@/stores/auth-token-store'
import { authenticationService } from '@/services/authentication-service'
import { storeToRefs } from 'pinia'

// TBD: This HTTP client currently calls the Federated Catalogue (FC) directly.
// In the future this may be replaced by Go backend endpoints that proxy
// FC requests.

const http = axios.create({
  // TODO: baseURL should be injected dynamically by ArgoCD / environment config
  // instead of being hard-coded here.
  baseURL: 'http://fc-server:8081',
  headers: { 'Content-Type': 'application/json' },
})

http.interceptors.request.use(
  (config) => {
    const tokenStore = useAuthTokenStore()
    const { isAuthSet, getAuthenticationHeader } = storeToRefs(tokenStore)
    if (isAuthSet.value) {
      config.headers.Authorization = getAuthenticationHeader.value
    } else if (config.headers && 'Authorization' in config.headers) {
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
      const res = await authenticationService.refresh()
      if (res) {
        return http(err.config)
      }
    }
    return Promise.reject(err)
  },
)

export default http
