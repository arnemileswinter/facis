import ContractTemplateListView from '@/views/contract-template-list/ContractTemplateListView.vue'
import LoginView from '@/views/auth/LoginView.vue'
import NewContractTemplateView from '@template-repository/views/NewContractTemplateView.vue'
import TableView from '@/views/TableView.vue'
import { DocumentTextIcon } from '@heroicons/vue/20/solid'
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import AuthSuccessView from '@/views/auth/AuthSuccessView.vue'
import LogoutCompleteView from '@/views/auth/LogoutCompleteView.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    meta: { name: 'DCS', hideInSidebar: true, requiresAuth: false, layout: 'blank', title: 'DCS' },
    component: LoginView,
  },
  {
    path: '/templates',
    name: 'templates.list',
    component: ContractTemplateListView,
    meta: { name: 'Contract Templates', icon: DocumentTextIcon, requiresAuth: true, title: 'DCS - Templates' },
  },
  {
    path: '/templates/new',
    name: 'templates.new',
    component: NewContractTemplateView,
    meta: { name: 'New Template', hideInSidebar: true, requiresAuth: true, title: 'DCS - New Template' },
  },
  {
    path: '/templates/edit/:did',
    name: 'templates.edit',
    component: NewContractTemplateView,
    meta: { name: 'Edit Template', hideInSidebar: true, requiresAuth: true, title: 'DCS - Edit Template' },
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
    meta: { name: 'Table', requiresAuth: false, title: 'DCS - Table' },
  },
  {
    path: '/auth/success',
    name: 'auth.success',
    meta: { hideInSidebar: true, requiresAuth: false, layout: 'blank', title: 'DCS - Auth Success' },
    component: AuthSuccessView,
  },
  {
    path: '/auth/logout-complete',
    name: 'auth.logout-complete',
    meta: { hideInSidebar: true, requiresAuth: false, layout: 'blank', title: 'DCS - Logout Complete' },
    component: LogoutCompleteView,
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes,
})

export { router }
