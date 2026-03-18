export interface SelfDescriptionsRequest {
  issuers?: string[]
  validators?: string[]
  statuses?: SelfDescriptionState[]
  ids?: string[]
  hashes?: string[]
  withMeta?: boolean
  withContent?: boolean
  offset?: number
  limit?: number
}

export interface FederatedCatalogueQueryRequest {
  statement: string
}

export type SelfDescriptionState = (typeof SelfDescriptionStateValue)[keyof typeof SelfDescriptionStateValue]
export const SelfDescriptionStateValue = {
  active: "ACTIVE",
  eol: "EOL",
  deprecated: "DEPRECATED",
  revoked: "REVOKED"
}