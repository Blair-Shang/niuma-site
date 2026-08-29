export type ProductStatus = 'ga' | 'beta' | 'coming-soon'

export interface ProductCardItem {
  id: string
  name: string
  tagline: string
  status: ProductStatus
  /** 产品详情路由；coming-soon 可为 null */
  to: string | null
  /** 真实产品界面截图；无图时卡片只展示文案 */
  shot: string | null
}

export interface ProductShot {
  src: string
  alt: string
  caption: string
  width: number
  height: number
}

/** 桌面端真实界面截图（首页 hero / 产品页共用） */
export const niumaDesktopShots: ProductShot[] = [
  {
    src: '/product/niuma-ops-workbench.png',
    alt: 'NiuMa 桌面端运维监控工作台：连接树、进程列表与 CPU 状态分析',
    caption: '运维监控工作台',
    width: 1024,
    height: 640,
  },
  {
    src: '/product/niuma-db-workbench.png',
    alt: 'NiuMa 桌面端数据库工作台：对象树、表数据浏览与日期编辑',
    caption: '数据库工作台',
    width: 1024,
    height: 640,
  },
]

/** 产品卡片封面：运维工作台 */
export const niumaDesktopShot = niumaDesktopShots[0].src

/** 工作室产品矩阵 — 后续新产品只追加此列表 */
export const products: ProductCardItem[] = [
  {
    id: 'niuma-desktop',
    name: 'NiuMa',
    tagline: '全能 AI 运维平台：数据库会话、SQL、监控与工具链一体。',
    status: 'ga',
    to: '/products/niuma',
    shot: niumaDesktopShot,
  },
  {
    id: 'placeholder-next',
    name: '更多产品',
    tagline: '后续开发中的工具将在此展示。',
    status: 'coming-soon',
    to: null,
    shot: null,
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
