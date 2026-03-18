import http from '@/api/catalogue-http'
import type { FederatedCatalogueQueryRequest, SelfDescriptionsRequest } from '@/models/requests/federated-catalogue-request';
import type { FederatedCatalogueQueryResponse, SelfDescriptionCreateResponse, SelfDescriptionDeleteResponse, SelfDescriptionsResponse } from '@/models/responses/federated-catalogue-response';
import type { SelfDescriptionContent } from '@/modules/template-catalogue/models/template-catalogue';

export const FederatedCatalogueService = {

  async getSelfDescriptions(request?: SelfDescriptionsRequest) {
    const params: Record<string, unknown> = { withContent: true }
    if (request?.issuers && request.issuers.length > 0) params.issuers = request.issuers
    if (request?.validators && request.validators.length > 0) params.validators = request.validators
    if (request?.statuses && request.statuses.length > 0) params.statuses = request.statuses
    if (request?.ids && request.ids.length > 0) params.ids = request.ids
    if (request?.hashes && request.hashes.length > 0) params.hashes = request.hashes
    if (request?.withMeta != null) params.withMeta = request.withMeta
    if (request?.withContent != null) params.withContent = request.withContent
    if (request?.offset != null) params.offset = request.offset
    if (request?.limit != null) params.limit = request.limit

    return http
      .get<SelfDescriptionsResponse>('/self-descriptions', { params })
      .then((res) => res.data)
      .catch((err) => {
        console.error('Get self-descriptions Error:', err)
        return {
          totalCount: 0,
          items: []
        } as SelfDescriptionsResponse
      })
  },

  async getSelfDescription(sdHash: string) {
    return http
      .get<SelfDescriptionContent>(`/self-descriptions/${encodeURIComponent(sdHash)}`)
      .then((res) => res.data)
      .catch((err) => {
        console.error('Get self-description error:', err.message)
        return null
      })
  },

  async createSelfDescription(payload: SelfDescriptionContent) {
    return http
      .post<SelfDescriptionCreateResponse>('/self-descriptions', payload)
      .then((res) => res.data)
  },

  async deleteSelfDescription(sdHash: string) {
    return http
      .delete<SelfDescriptionDeleteResponse>(`/self-descriptions/${encodeURIComponent(sdHash)}`)
      .then((res) => res.data)
  },

  async query<T = unknown>(statement: string) {
    const payload: FederatedCatalogueQueryRequest = { statement }
    return http
      .post<FederatedCatalogueQueryResponse<T>>('/query', payload)
      .then((res) => res.data)
      .catch((err) => {
        console.error('Query Error:', err)
        return { totalCount: 0, items: [] } as FederatedCatalogueQueryResponse<T>
      })
  },
}

