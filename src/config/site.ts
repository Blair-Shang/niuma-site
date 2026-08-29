/** 站点公开配置（构建期环境变量） */

export const siteConfig = {
  name: 'NiuMa',
  tagline: '专业开发与运维工具',
  description:
    'NiuMa — 工作室品牌与旗舰 AI 运维桌面平台。数据库、SQL、监控与工具链一体。',
  apiBase: (import.meta.env.VITE_SITE_API_BASE || '').replace(/\/$/, ''),
  /**
   * niuma-cloud API 基址（含 /niuma/cloud 前缀，不含 /api/v1）。
   * 开发默认本地；生产默认 https://www.niuma007.com/niuma/cloud
   */
  cloudApiBase: (
    import.meta.env.VITE_CLOUD_API_BASE ||
    (import.meta.env.DEV ? 'http://127.0.0.1:8090/niuma/cloud' : 'https://www.niuma007.com/niuma/cloud')
  ).replace(/\/$/, ''),
  /**
   * 下载页版本/说明回落（主数据源为 cloud updates/latest）。
   * 真实安装包由服务端 hit 302 到 cloud /updates/download。
   * Windows=stable；Linux / macOS=beta。
   */
  download: {
    version: import.meta.env.VITE_DOWNLOAD_VERSION || '',
    windowsUrl: import.meta.env.VITE_DOWNLOAD_WIN_URL || '',
  },
  /**
   * ICP 网站备案号（首页底部悬挂，链到工信部备案系统）。
   * 湖北非广东，应挂服务/网站号（主体号后带序号），不要只挂主体号。
   */
  icp: {
    number: '鄂ICP备2026041693号-1',
    url: 'https://beian.miit.gov.cn/',
  },
} as const

export function apiUrl(path: string): string {
  const p = path.startsWith('/') ? path : `/${path}`
  return `${siteConfig.apiBase}${p}`
}
