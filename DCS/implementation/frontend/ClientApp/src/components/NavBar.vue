<template>
  <div class="drawer lg:drawer-open">
    <input id="main-drawer" type="checkbox" class="drawer-toggle" />
    
    <div class="drawer-content flex flex-col min-h-screen bg-base-100">
      <header class="navbar w-full bg-base-200 border-b border-base-content/10 sticky top-0 z-30">
        <div class="flex-none">
          <label for="main-drawer" class="btn btn-square btn-ghost lg:hidden">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="w-6 h-6 stroke-current">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
            </svg>
          </label>
          
          <button 
            @click="isCollapsed = !isCollapsed" 
            class="btn btn-square btn-ghost hidden lg:inline-flex"
          >
            <svg 
              xmlns="http://www.w3.org/2000/svg" 
              fill="none" 
              viewBox="0 0 24 24" 
              stroke-width="1.5" 
              stroke="currentColor" 
              :class="['w-6 h-6 transition-transform duration-300', isCollapsed ? 'rotate-180' : '']"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d="M18.75 19.5l-7.5-7.5 7.5-7.5m-6 15L5.25 12l7.5-7.5" />
            </svg>
          </button>
        </div>
        
        <div class="flex-1 px-4 font-bold text-lg">
          DCS <span class="font-normal opacity-50 ml-2">| Dashboard</span>
        </div>
      </header>

      <main class="flex-grow p-4 md:p-8">
        <RouterView />
      </main>
    </div>

    <div class="drawer-side z-40">
      <label for="main-drawer" aria-label="close sidebar" class="drawer-overlay"></label>
      
      <aside 
        :class="[
          'flex flex-col min-h-full bg-base-200 border-r border-base-content/5 transition-all duration-300 ease-in-out',
          isCollapsed ? 'lg:w-20' : 'w-72'
        ]"
      >
        <div class="flex items-center h-16 px-4 overflow-hidden">
          <div  class="font-bold text-2xl tracking-tight text-base-content uppercase">DCS</div>
        </div>

        <nav class="flex-1 overflow-y-auto overflow-x-hidden py-4">
          <ul class="menu px-3 gap-1 w-full text-base-content">
            <li v-for="route in navigationRoutes" :key="route.path">
              <RouterLink 
                :to="route.path"
                @click="closeMobileDrawer"
                :class="[
                  'flex items-center gap-4 py-3 rounded-btn',
                  isCollapsed ? 'justify-center px-0' : 'px-4'
                ]"
                active-class="active bg-primary text-primary-content"
                :data-tip="isCollapsed ? route.meta.name : ''"
              >
                <component 
                  :is="route.meta?.icon" 
                  class="w-6 h-6 flex-shrink-0" 
                  aria-hidden="true" 
                />
                <span v-if="!isCollapsed" class="font-medium whitespace-nowrap">
                  {{ route.meta.name }}
                </span>
              </RouterLink>
            </li>
          </ul>
        </nav>

        <div class="p-4 border-t border-base-content/10 bg-base-300/20">
          <div :class="['flex items-center gap-3', isCollapsed ? 'justify-center' : 'px-2']">
            <div class="avatar">
              <div class="w-10 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
                <img src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=facearea&facepad=2&w=128&h=128&q=80" alt="Profile" />
              </div>
            </div>
            <div v-if="!isCollapsed" class="overflow-hidden">
              <p class="text-sm font-bold truncate">Tom Cook</p>
              <p class="text-xs opacity-60">Admin</p>
            </div>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, RouterLink, RouterView } from 'vue-router'

const router = useRouter()
const isCollapsed = ref(false)

const closeMobileDrawer = () => {
  const drawerToggle = document.getElementById('main-drawer') as HTMLInputElement | null
  if (drawerToggle) drawerToggle.checked = false
}

const navigationRoutes = computed(() => {
  try {
    return router.getRoutes().filter(route =>
      route.name &&
      !route.path.includes(':') &&
      route.meta?.name &&
      route.meta?.hideInSidebar !== true
    )
  } catch (e) {
    return []
  }
})
</script>