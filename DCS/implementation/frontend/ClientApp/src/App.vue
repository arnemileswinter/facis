<script setup lang="ts">
import NavBar from './components/NavBar.vue';

</script>

<template>
  <header>
    <div class="bg-blue-500 text-white p-4 rounded">
      Hallo Tailwind 4 + Vue 3!
    </div>
    <NavBar />
  </header>
  <main>
    <RouterView />
  </main>
</template>