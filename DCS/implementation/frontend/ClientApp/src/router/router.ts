import ApproveContractTemplateView from '@/modules/template-repository/views/ApproveContractTemplateView.vue'
import ViewContractTemplateView from '@/modules/template-repository/views/ViewContractTemplateView.vue'
import { useAuthStore } from '@/stores/auth-store'
import AuthSuccessView from '@/views/auth/AuthSuccessView.vue'
import LoginView from '@/views/auth/LoginView.vue'
import ContractTemplateListView from '@/views/contract-template-list/ContractTemplateListView.vue'
import { AuthenticationService } from '@/services/authentication-service'
import { DocumentCheckIcon, DocumentMagnifyingGlassIcon, DocumentTextIcon } from '@heroicons/vue/20/solid'
import NewContractTemplateView from '@template-repository/views/NewContractTemplateView.vue'
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { getUIBasePath } from '@/config'
import ReviewContractTemplateView from '@/modules/template-repository/views/ReviewContractTemplateView.vue'
import ContractTemplateTaskView from '@/views/contract-template-list/ContractTemplateTaskView.vue'

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
    meta: { name: 'Contract Templates', icon: DocumentTextIcon, requiresAuth: true, title: 'DCS - Templates', order: 1 },
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
  },
  {
    path: '/templates/view/:did',
    name: 'templates.view',
    component: ViewContractTemplateView,
    meta: { name: 'View Template', hideInSidebar: true, requiresAuth: true, title: 'DCS - View Template' },
  },
  {
    path: '/templates/review/:did',
    name: 'templates.review',
    component: ReviewContractTemplateView,
    meta: { name: 'Review Template', hideInSidebar: true, requiresAuth: true, title: 'DCS - Review Template' },
  },
  {
    path: '/templates/approve/:did',
    name: 'templates.approve',
    component: ApproveContractTemplateView,
    meta: { name: 'Approve Template', hideInSidebar: true, requiresAuth: true, title: 'DCS - Approve Template' },
  },
  {
    path: '/templates/tasks/review',
    name: 'templates.tasks.review',
    component: ContractTemplateTaskView,
    meta: {
      name: 'Assigned Review Tasks',
      icon: DocumentMagnifyingGlassIcon,
      requiresAuth: true,
      title: 'DCS - Review Tasks',
      order: 2,
      roles: ['TEMPLATE_REVIEWER']
    },
  },
  {
    path: '/templates/tasks/approval',
    name: 'templates.tasks.approval',
    component: ContractTemplateTaskView,
    meta: {
      name: 'Assigned Approval Tasks',
      icon: DocumentCheckIcon,
      requiresAuth: true,
      title: 'DCS - Approval Tasks',
      order: 3,
      roles: ['TEMPLATE_APPROVER']
    },
  },
  {
    path: '/auth/success',
    name: 'auth.success',
    meta: { hideInSidebar: true, requiresAuth: false, layout: 'blank', title: 'DCS - Auth Success' },
    component: AuthSuccessView,
  },
]

const router = createRouter({
  history: createWebHistory(getUIBasePath()),
  routes: routes,
})

router.beforeEach(async (to) => {
  if (to.meta.requiresAuth === false) {
    return true
  }

  const authStore = useAuthStore()
  if (authStore.isAuthenticated) {
    return true
  }

  await AuthenticationService.refresh()
  if (authStore.isAuthenticated) {
    return true
  }

  const loginUrl = await AuthenticationService.getLoginPath()
  if (loginUrl) {
    window.location.href = loginUrl
    return false
  }

  return { name: 'home' }
})

export { router }
