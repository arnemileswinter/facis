import type { FunctionalComponent, HTMLAttributes, VNodeProps } from 'vue'
import 'vue-router'

export {}

declare module 'vue-router' {
  interface RouteMeta {
    name?: string
    hideInSidebar?: boolean
    icon?: FunctionalComponent<HTMLAttributes & VNodeProps>
  }
}
