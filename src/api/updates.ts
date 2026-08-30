import { siteConfig } from '../config/site'

/** cloud 公开最新发布（与桌面端共用） */
export type UpdateRelease = {
  product: string
  channel: string
  platform: string
  arch: string
  version: string
  title: string
  notesMd: string
  downloadUrl: string
  sha256: string
  fileSize: number
  publishedAt?: string
}

export type FetchLatestParams = {
  platform?: string
  arch?: string
  product?: string
  channel?: string
}

function cloudURL(path: string): string {
  const base = siteConfig.cloudApiBase.replace(/\/$/, '')
  const p = path.startsWith('/') ? path : `/${path}`
  return `${base}${p}`
}

function releaseQuery(params: FetchLatestParams = {}): URLSearchParams {
  const q = new URLSearchParams()
  q.set('platform', params.platform || 'windows')
  q.set('arch', params.arch || 'x64')
  if (params.product) q.set('product', params.product)
  if (params.channel) q.set('channel', params.channel)
  return q
}

/** 拉取 cloud 最新 published 元数据（展示版本与说明）。安装包请走官网 hit，勿用本结果的 downloadUrl 直接跳转。 */
export async function fetchLatestRelease(
  params: FetchLatestParams = {},
): Promise<UpdateRelease | null> {
  try {
    const url = cloudURL(`/api/v1/updates/latest?${releaseQuery(params)}`)
    const res = await fetch(url, { headers: { Accept: 'application/json' } })
    if (res.status === 404) return null
    if (!res.ok) return null
    return (await res.json()) as UpdateRelease
  } catch {
    return null
  }
}

/** 已发布版本列表（semver 降序），供下载页更新说明。 */
export async function fetchReleaseHistory(
  params: FetchLatestParams & { limit?: number } = {},
): Promise<UpdateRelease[]> {
  const q = releaseQuery(params)
  if (params.limit && params.limit > 0) {
    q.set('limit', String(params.limit))
  }
  try {
    const url = cloudURL(`/api/v1/updates/releases?${q}`)
    const res = await fetch(url, { headers: { Accept: 'application/json' } })
    if (!res.ok) return []
    const data = (await res.json()) as { items?: UpdateRelease[] }
    return Array.isArray(data.items) ? data.items : []
  } catch {
    return []
  }
}
