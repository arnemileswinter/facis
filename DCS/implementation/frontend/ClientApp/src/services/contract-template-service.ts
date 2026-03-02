import axios from 'axios'
import type { ContractTemplate } from '../models/contract-template'
import type {
  ContractTemplateApproveRequest,
  ContractTemplateCreateRequest,
  ContractTemplateRejectRequest,
  ContractTemplateRetrieveByIdRequest,
  ContractTemplateRetrieveRequest,
  ContractTemplateSearchRequest,
  ContractTemplateSubmitRequest,
  ContractTemplateUpdateRequest,
} from '../models/requests/template-request'
import type {
  ContractTemplateApproveResponse,
  ContractTemplateCreateResponse,
  ContractTemplateRejectResponse,
  ContractTemplateRetrieveByIdResponse,
  ContractTemplateRetrieveResponse,
  ContractTemplateSearchResponse,
  ContractTemplateSubmitResponse,
  ContractTemplateUpdateResponse,
} from '../models/responses/template-response'
import { AuthenticationService } from './authentication-service'

const API_BASE_URL = import.meta.env.DCS_API_BASE_URL

const token_type = localStorage.getItem('token_type')
const access_token = localStorage.getItem('access_token')

const http = axios.create({
  baseURL: API_BASE_URL,
  headers: { 'Content-Type': 'application/json' },
})

if (token_type && access_token) {
  http.defaults.headers.common.Authorization = `${token_type} ${access_token}`
}

http.interceptors.request.use((config) => {
  const token_type = localStorage.getItem('token_type')
  const access_token = localStorage.getItem('access_token')
  if (token_type && access_token) {
    config.headers.setAuthorization(`${token_type} ${access_token}`)
  }
  return config
})

http.interceptors.response.use(
  (resp) => resp,
  (err) => {
    console.log('Reject:', err)
    AuthenticationService.refresh()
  },
)

export const ContractTemplateService = {
  async create(request: ContractTemplateCreateRequest) {
    return http.post<ContractTemplateCreateResponse>('/template/create', request).then((res) => res.data)
  },

  async submit(request: ContractTemplateSubmitRequest) {
    return http.post<ContractTemplateSubmitResponse>('/template/submit', request).then((res) => res.data)
  },

  async update(request: ContractTemplateUpdateRequest) {
    return http.put<ContractTemplateUpdateResponse>('/template/update', request).then((res) => res.data)
  },

  async search(_request: ContractTemplateSearchRequest) {
    return http.get<ContractTemplateSearchResponse>('template/search').then((res) => res.data)
  },

  async retrieve(_request?: ContractTemplateRetrieveRequest): Promise<ContractTemplate[]> {
    return http
      .get<ContractTemplateRetrieveResponse>('/template/retrieve')
      .then((res) => {
        return Array.isArray(res.data.contract_templates) ? res.data.contract_templates : []
      })
      .catch((err) => {
        console.error('Retrieve Error:', err)
        return []
      })
  },

  async retrieveById(request: ContractTemplateRetrieveByIdRequest): Promise<ContractTemplate | null> {
    const queryParams = { document_number: request.document_number, version: request.version }
    return http
      .get<ContractTemplateRetrieveByIdResponse>(`/template/retrieve/${request.did}`, { params: queryParams })
      .then((res) => {
        console.log(res.status)
        return { ...res.data }
      })
      .catch((err) => {
        console.error('Retrieve ID Error:', err.message)
        return null
      })
  },

  async approve(request: ContractTemplateApproveRequest) {
    return http.post<ContractTemplateApproveResponse>('/template/approve', request).then((res) => res.data)
  },

  async reject(request: ContractTemplateRejectRequest) {
    return http.post<ContractTemplateRejectResponse>('/template/reject', request).then((res) => res.data)
  },
}
