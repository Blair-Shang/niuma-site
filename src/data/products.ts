export type ProductStatus = 'ga' | 'beta' | 'coming-soon'

export interface ProductCardItem {
  id: string
  name: string
  tagline: string
  status: ProductStatus
  /** 产品详情路由；coming-soon 可为 null */
  to: string | null
}

/** 工作室产品矩阵 — 后续新产品只追加此列表 */
export const products: ProductCardItem[] = [
  {
    id: 'niuma-desktop',
    name: 'NiuMa',
    tagline: '全能 AI 运维平台：数据库会话、SQL、监控与工具链一体。',
    status: 'ga',
    to: '/products/niuma',
  },
  {
    id: 'placeholder-next',
    name: '更多产品',
    tagline: '后续开发中的工具将在此展示。',
    status: 'coming-soon',
    to: null,
  },
]

export const statusLabel: Record<ProductStatus, string> = {
  ga: '正式版',
  beta: 'Beta',
  'coming-soon': '即将推出',
}

export const statusBadgeVariant: Record<
  ProductStatus,
  'success' | 'warning' | 'default'
> = {
  ga: 'success',
  beta: 'warning',
  'coming-soon': 'default',
}
