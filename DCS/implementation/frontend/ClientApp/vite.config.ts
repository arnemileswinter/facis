import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { fileURLToPath } from 'url'

// https://vite.dev/config/
export default defineConfig(({mode}) => {
  const env = loadEnv(mode, process.cwd(), '')
  return {
    plugins: [vue(), tailwindcss()],
    envPrefix: 'DCS',
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src/', import.meta.url)),
        '@core': fileURLToPath(new URL('./src/core/', import.meta.url)),
        '@template-repository': fileURLToPath(new URL('./src/modules/template-repository/', import.meta.url)),
      },
          proxy: {
      '/api': {
        target: "http://localhost:8991",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    },
    server: {
      proxy: {
        '/api': {
          target: env.DCS_API_BASE_URL,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ''),
        },
      },
    },

  }
})
