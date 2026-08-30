import { apiUrl, siteConfig } from '../config/site'

export type DownloadPlatform = 'windows' | 'linux' | 'macos'

export type DownloadStats = {
  total: number
  byPlatform?: Record<string, number>
}

/** 拉取全局累计下载（cloud：官网 + 桌面自动更新）；失败返回 null，UI 隐藏数字 */
export async function fetchDownloadStats(): Promise<DownloadStats | null> {
  const base = siteConfig.cloudApiBase.replace(/\/$/, '')
  try {
    const res = await fetch(`${base}/api/v1/updates/hits?product=niuma`, {
      headers: { Accept: 'application/json' },
    })
    if (!res.ok) return null
    const data = (await res.json()) as DownloadStats
    if (typeof data.total !== 'number' || Number.isNaN(data.total)) return null
    return data
  } catch {
    return null
  }
}

/**
 * 触发安装包下载：走同仓 Go API，302 到 cloud 最新包（由 cloud 记全局次数）。
 * Windows 为 stable；Linux / macOS 为 beta。
 * 仅当 VITE_DOWNLOAD_USE_DIRECT=true 且为 Windows 时直链回退（不跟随流水线新版本）。
 */
export function startPlatformDownload(platform: DownloadPlatform): void {
  if (
    platform === 'windows' &&
    import.meta.env.VITE_DOWNLOAD_USE_DIRECT === 'true' &&
    siteConfig.download.windowsUrl
  ) {
    window.open(siteConfig.download.windowsUrl, '_blank', 'noopener,noreferrer')
    return
  }
  window.location.assign(apiUrl(`/api/v1/downloads/${platform}/hit`))
}
