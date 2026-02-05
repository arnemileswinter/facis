<script setup lang="ts">
import NavBar from './components/NavBar.vue'

</script>
<template>
  <div class="min-h-screen bg-gray-800">
    <header>
      <NavBar />
    </header>
    <main class="md:ml-64 min-h-screen">
      <div class="p-6 lg:p-10">
        <div class="max-w-7xl mx-auto">
          <RouterView />
        </div>
      </div>
    </main>
  </div>
</template>
