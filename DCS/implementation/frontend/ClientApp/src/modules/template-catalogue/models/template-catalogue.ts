import type { TemplateTypeValue } from "@/modules/template-repository/models/contract-templace"

export interface FederatedCatalogueSdMeta {
  /** credentialSubject ID */
  id: string
  sdHash: string
  issuer: string
  uploadDatetime: string
  statusDatetime: string
}

export interface TemplateCatalogue {
  did: string
  documentNumber?: string
  version?: string
  description?: string
  name?: string
  templateType?: TemplateTypeValue | string
  participantId?: string
  createdAt?: string
  updatedAt?: string
  sdMeta?: FederatedCatalogueSdMeta
}
export interface SelfDescriptionContent { [key: string]: any }
