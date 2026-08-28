import { apiUrl, siteConfig } from '../config/site'

export type DownloadStats = {
  total: number
  byPlatform?: Record<string, number>
}

/** 拉取累计下载；失败返回 null，UI 隐藏数字 */
export async function fetchDownloadStats(): Promise<DownloadStats | null> {
  try {
    const res = await fetch(apiUrl('/api/v1/downloads/stats'), {
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
 * 触发 Windows 下载：默认走同仓 Go API（记数后 302）。
 * 仅当 VITE_DOWNLOAD_USE_DIRECT=true 时直链回退（无后端预览）。
 */
export function startWindowsDownload(): void {
  if (
    import.meta.env.VITE_DOWNLOAD_USE_DIRECT === 'true' &&
    siteConfig.download.windowsUrl
  ) {
    window.open(siteConfig.download.windowsUrl, '_blank', 'noopener,noreferrer')
    return
  }
  window.location.assign(apiUrl('/api/v1/downloads/windows/hit'))
}
