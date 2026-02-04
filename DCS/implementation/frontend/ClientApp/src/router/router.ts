import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import TableView from '../views/TableView.vue'
import ContractTemplateView from '../views/contract-template/ContractTemplateView.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: ContractTemplateView,
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
