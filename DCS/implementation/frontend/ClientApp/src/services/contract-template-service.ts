import axios from 'axios'
import type { ContractTemplate } from '../models/contract-template'
import type {
  ContractTemplateRetrieveByIdRequest,
  ContractTemplateSearchRequest,
} from '../models/requests/template-request'
import type {
  ContractTemplateRetrieveByIdResponse,
  ContractTemplateRetrieveResponse,
} from '../models/responses/template-response'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL

const http = axios.create({
  baseURL: API_BASE_URL,
  headers: { 'Content-Type': 'application/json' },
})

export const ContractTemplateService = {
  async create(data: any) {
    return http.post('/template/create', data).then((res) => res.data)
  },

  async submit(data: any) {
    return http.post('/template/submit', data).then((res) => res.data)
  },

  async update(data: string) {
    return http.put('/template/update', data).then((res) => res.data)
  },

  async retrieve(): Promise<ContractTemplate[]> {
    return http
      .get<ContractTemplateRetrieveResponse>('/template/retrieve')
      .then((res) => {
        return Array.isArray(res.data.contract_templates)
          ? res.data.contract_templates
          : []
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

  async search(_request: ContractTemplateSearchRequest) {
    return http.get('template/search').then((res) => res.data)
  },

  async approve(data: any) {
    return http.post('/template/approve', data).then((res) => res.data)
  },

  async reject(data: any) {
    return http.post('/template/reject', data).then((res) => res.data)
  },
}
