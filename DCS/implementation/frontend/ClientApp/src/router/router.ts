import { useAuthStore } from '@/stores/auth-store'
import ContractTemplateListView from '@/views/contract-template-list/ContractTemplateListView.vue'
import LoginView from '@/views/login/LoginView.vue'
import NewContractTemplateView from '@/views/new-contract-template/NewContractTemplateView.vue'
import TableView from '@/views/TableView.vue'
import { DocumentTextIcon } from '@heroicons/vue/20/solid'
import { storeToRefs } from 'pinia'
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    meta: { name: 'Home', hideInSidebar: true, requiresAuth: false, title: 'DCS' },
    redirect: { name: 'templates.list' },
  },
  {
    path: '/templates',
    name: 'templates.list',
    component: ContractTemplateListView,
    meta: { name: 'Contract Templates', icon: DocumentTextIcon, requiresAuth: false, title: 'DCS - Templates' },
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView,
    meta: { name: 'Login', hideInSidebar: true, requiresAuth: false, title: 'DCS - Login' },
    beforeEnter: () => {
      const authStore = useAuthStore()
      const { isAuthenticated } = storeToRefs(authStore)
      console.log(isAuthenticated.value)
      return !isAuthenticated.value
    },
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
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes,
})

router.beforeEach((to, _from) => {
  const authStore = useAuthStore()
  const { isAuthenticated } = storeToRefs(authStore)
  if (to.meta.requiresAuth && !isAuthenticated.value) {
    return { name: 'login' }
  }
})

export { router }
