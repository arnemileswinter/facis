<template>
  <nav class="bg-blue-300 p-4">
    <div>
      <template v-for="route in $router.getRoutes()" :key="route.path">
        <RouterLink :to="route.path" class="px-4">{{ route.name }}</RouterLink>
      </template>
    </div>
  </nav>
</template>
