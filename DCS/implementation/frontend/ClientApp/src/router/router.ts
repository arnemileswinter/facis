import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import TableView from '../views/TableView.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: HomeView,
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
