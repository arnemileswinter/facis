import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { fileURLToPath } from 'url'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(),
  tailwindcss(),
  ],
  envPrefix: 'DCS',
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src/', import.meta.url)),
      '@core': fileURLToPath(new URL('./src/core/', import.meta.url)),
      '@template-repository': fileURLToPath(new URL('./src/modules/template-repository/', import.meta.url)),
    }
  }
})
