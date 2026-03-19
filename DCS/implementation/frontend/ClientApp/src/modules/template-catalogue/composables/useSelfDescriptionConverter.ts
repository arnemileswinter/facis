import { ref } from "vue";
import type { FederatedCatalogueSdMeta, TemplateCatalogue } from "../models/template-catalogue";
import { FederatedCatalogueService } from "@/services/federated-catalogue-services";
import { SelfDescriptionStateValue } from "@/models/requests/federated-catalogue-request";

async function enrichWithSdMeta(
  catalogues: Omit<TemplateCatalogue, "sdMeta">[], didById: Record<string, string>
): Promise<TemplateCatalogue[]> {
  const metaByDid: Record<string, FederatedCatalogueSdMeta> = {}
  const metaResp = await FederatedCatalogueService.getSelfDescriptions({
    ids: Object.keys(didById),
    withMeta: true,
    withContent: false,
    statuses: [SelfDescriptionStateValue.active],
  })

  metaResp.items.forEach((item) => {
    const id = item?.meta?.id
    const sdHash = item?.meta?.sdHash
    const issuer = item?.meta?.issuer
    const uploadDatetime = item?.meta?.uploadDatetime
    const statusDatetime = item?.meta?.statusDatetime
    const did = didById[id]

    if (!did || !id || !sdHash || !issuer || !uploadDatetime || !statusDatetime) return

    metaByDid[did] = { id, sdHash, issuer, uploadDatetime, statusDatetime }
  })

  return catalogues.map((c) => ({
    ...c,
    sdMeta: metaByDid[c.did],
  }))
}

export function useSelfDescriptionConverter() {
  const loading = ref(false)
  const error = ref<string | null>(null)
  const totalCount = ref(0)

  async function loadTemplateCatalogues(page: number, pageSize: number): Promise<TemplateCatalogue[]> {
    loading.value = true
    error.value = null
    let result: TemplateCatalogue[] = []
    try {
      const skip = page * pageSize
      const limit = pageSize

      const statement = [
        "MATCH (n:ContractTemplate) ",
        "RETURN n, n.claimsGraphUri AS claimsGraphUri ",
        "ORDER BY n.createdAt DESC ",
        `SKIP ${skip} LIMIT ${limit}`,
      ].join("")

      const queryResp = await FederatedCatalogueService.query<{ n: any; claimsGraphUri?: string[] }>(statement)
      totalCount.value = queryResp.totalCount ?? queryResp.items.length

      const catalogues: Omit<TemplateCatalogue, "sdMeta">[] = []
      const didById: Record<string, string> = {}


      queryResp.items.forEach((row) => {
        const node = row?.n
        if (!node) return
        const id = node.claimsGraphUri?.[0]
        if (!id) return

        catalogues.push({
          did: node.did,
          documentNumber: node.documentNumber,
          version: node.version,
          name: node.name,
          description: node.description,
          templateType: node.templateType,
          participantId: node.participantId,
          createdAt: node.createdAt,
          updatedAt: node.updatedAt,
        })
        didById[id] = node.did
      })

      result = await enrichWithSdMeta(catalogues, didById)
    } catch (err: any) {
      error.value = err.message || "Error loading template catalogues"
    } finally {
      loading.value = false
    }
    return result
  }

  async function loadTemplateCatalogue(did: string): Promise<TemplateCatalogue | null> {
    loading.value = true
    error.value = null
    let result: TemplateCatalogue | null = null
    try {
      if (!did) return result

      const statement = [
        "MATCH (n:ContractTemplate) ",
        `WHERE (n.did = '${did}') `,
        "RETURN n, n.claimsGraphUri AS claimsGraphUri ",
        "LIMIT 1",
      ].join("")

      const queryResp = await FederatedCatalogueService.query<{ n: any; claimsGraphUri?: string[] }>(statement)
      const row = queryResp.items?.[0]
      const node = (row as any)?.n ?? (row as any)?.["n"]
      const didById: Record<string, string> = {}
      if (!node) return result

      const baseCatalogue: Omit<TemplateCatalogue, "sdMeta"> = {
        did: node.did ?? did,
        documentNumber: node.documentNumber,
        version: node.version,
        name: node.name,
        description: node.description,
        templateType: node.templateType,
        participantId: node.participantId,
        createdAt: node.createdAt,
        updatedAt: node.updatedAt,
      }
      const id = row?.claimsGraphUri?.[0] ?? ""
      didById[id] = node.did
      const catalogues = await enrichWithSdMeta([baseCatalogue], didById)
      if (catalogues.length > 0) result = catalogues[0] ?? null
      else result = baseCatalogue
    } catch (err: any) {
      error.value = err.message || "Error loading template catalogue"
    } finally {
      loading.value = false
    }

    return result
  }


  return {
    loading,
    error,
    totalCount,
    loadTemplateCatalogues,
    loadTemplateCatalogue
  }
}