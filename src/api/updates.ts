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

/** 拉取 cloud 最新 published；无发布或失败返回 null（UI 回落 env）。 */
export async function fetchLatestRelease(
  params: FetchLatestParams = {},
): Promise<UpdateRelease | null> {
  const q = new URLSearchParams()
  q.set('platform', params.platform || 'windows')
  q.set('arch', params.arch || 'x64')
  if (params.product) q.set('product', params.product)
  if (params.channel) q.set('channel', params.channel)
  try {
    const res = await fetch(cloudURL(`/api/v1/updates/latest?${q}`), {
      headers: { Accept: 'application/json' },
    })
    if (res.status === 404) return null
    if (!res.ok) return null
    return (await res.json()) as UpdateRelease
  } catch {
    return null
  }
}
