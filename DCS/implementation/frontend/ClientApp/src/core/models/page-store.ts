interface PageState {
  title: string
  isSidebarCollapsed: boolean
  breadcrumbs?: BreadcrumbItem[]
  pageSidebarId: string
}

interface BreadcrumbItem {
  label: string
  path?: string
}

export type { PageState, BreadcrumbItem }