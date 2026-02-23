import { DocumentTextIcon } from '@heroicons/vue/20/solid'
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import TableView from '../views/TableView.vue'
import ContractTemplateView from '../views/contract-template-list/ContractTemplateListView.vue'
import NewContractTemplateView from '../views/new-contract-template/NewContractTemplateView.vue'
import { DocumentTextIcon } from '@heroicons/vue/20/solid'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'templates.list',
    component: ContractTemplateView,
    meta: { name: 'Contract Templates', icon: DocumentTextIcon },
  },
  {
    path: '/templates/new',
    name: 'templates.new',
    component: NewContractTemplateView,
    meta: { name: 'New Template', hideInSidebar: true },
  },
  {
    path: '/templates/edit/:did',
    name: 'templates.edit',
    component: NewContractTemplateView,
    meta: { name: 'Edit Template', hideInSidebar: true },
    props: (route) => {
      const did = route.params.did
      const document_number = route.query.document_number
      const version = route.query.version
      if (
        did &&
        document_number &&
        version &&
        !Array.isArray(did) &&
        !Array.isArray(document_number) &&
        !Array.isArray(version)
      ) {
        return {
          did: did,
          document_number: parseInt(document_number),
          version: parseInt(version),
        }
      }
    },
  },
  {
    path: '/table',
    name: 'table',
    component: TableView,
    meta: { name: 'Table' },
  },
]

declare module 'vue-router' {
  interface RouteMeta {
    title: string
    icon?: unknown
    hideInSidebar?: boolean
  }
}

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes,
})
