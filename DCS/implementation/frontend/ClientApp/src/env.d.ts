/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<object, object, unknown>
  export default component
}

interface ViteTypeOptions {
  strictImportMetaEnv: unknown
}

interface ImportMetaEnv {
  readonly DCS_API_BASE_URL: string
  readonly DCS_API_URL: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}