import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import TableView from '../views/TableView.vue'
import ContractTemplateView from '../views/contract-template/ContractTemplateView.vue'
import { DocumentTextIcon } from '@heroicons/vue/20/solid'

const routes = [
  {
    path: '/',
    name: 'Contract Templates',
    component: ContractTemplateView,
    meta: { icon: DocumentTextIcon }
  },
  {
    path: '/table',
    name: 'Table',
    component: TableView,
  },
]

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes,
})
