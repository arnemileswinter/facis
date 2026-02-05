<template>
  <div class="min-h-screen bg-white">
    <div class="hidden md:fixed md:inset-y-0 md:flex md:w-72 md:flex-col">
      <div class="flex flex-col flex-grow bg-gray-900 pt-5 overflow-y-auto border-r border-white/5">

        <div class="flex items-center flex-shrink-0 px-6 mb-8">
          <!-- <img class="h-8 w-auto" src="https://tailwindui.com/plus/img/logos/mark.svg?color=indigo&shade=500" alt="Logo" /> -->
          <div class=" mb-10 text-white font-bold text-2xl tracking-tight">
            DCS </div>
        </div>
        <nav class="flex-1 px-3 space-y-1">
          <RouterLink v-for="route in navigationRoutes" :key="route.path" :to="route.path" v-slot="{ isActive }"
            class="group flex items-center px-3 py-2 text-sm font-semibold rounded-md transition-colors duration-200"
            :class="[
              isActive
                ? 'bg-gray-800 text-white'
                : 'text-gray-400 hover:text-white hover:bg-gray-800'
            ]">
            <component :is="route.meta?.icon" class="flex-shrink-0 w-6 h-6 mr-3 transition-colors duration-200"
              :class="[isActive ? 'text-white' : 'text-gray-400 group-hover:text-white']" aria-hidden="true" />

            {{ route.name }}
          </RouterLink>
        </nav>

        <div class="flex-shrink-0 flex p-4 bg-gray-900">
          <div class="flex items-center w-full">
            <img class="h-9 w-9 rounded-full ring-2 ring-white/10"
              src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
              alt="" />
            <div class="ml-3">
              <p class="text-sm font-medium text-white">Tom Cook</p>
              <p class="text-xs font-medium text-gray-400">View profile</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="md:pl-72">
      <main class="py-10 px-8">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const navigationRoutes = computed(() => {
  try {
    return router.getRoutes().filter(route =>
      route.name &&
      !route.path.includes(':') &&
      route.meta?.hideInSidebar !== true
    )
  } catch (e) {
    return []
  }
})
</script>