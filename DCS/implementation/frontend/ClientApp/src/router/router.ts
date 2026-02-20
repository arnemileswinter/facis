import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import TableView from '../views/TableView.vue'
import ContractTemplateView from '../views/contract-template-list/ContractTemplateListView.vue'
import NewContractTemplateView from '../views/new-contract-template/NewContractTemplateView.vue'
import { DocumentTextIcon } from '@heroicons/vue/20/solid'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Contract Templates',
    component: ContractTemplateView,
    meta: { icon: DocumentTextIcon, title: "DCS" }
  },
  {
    path: '/templates/new',
    name: 'New Template',
    component: NewContractTemplateView,
    meta: { hideInSidebar: true, title: "DCS - New Template" }
  },
  {
    path: '/table',
    name: 'Table',
    component: TableView,
    meta: { title: "DCS - Table" }
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
