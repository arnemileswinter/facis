import type { SelfDescriptionState } from "@/models/requests/federated-catalogue-request"

export interface SelfDescriptionsResponse {
  totalCount: number,
  items: {
    meta: {
      expirationTime: null,
      content: string | null,
      validators: null,
      sdHash: string,
      id: string,
      status: SelfDescriptionState,
      issuer: string,
      validatorDids: null,
      uploadDatetime: string,
      statusDatetime: string
    },
    content: string | null
  }[]
}

export interface SelfDescriptionCreateResponse {
  sdHash: string
  id: string
  status: string
  issuer: string
  validatorDids: null | string,
  uploadDatetime: string
  statusDatetime: string
}

export interface SelfDescriptionDeleteResponse {
  code?: string,
  message?: string
}

export interface FederatedCatalogueQueryResponse<T = unknown> {
  totalCount?: number
  items: T[]
}