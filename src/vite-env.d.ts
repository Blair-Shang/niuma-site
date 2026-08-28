/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_SITE_API_BASE: string
  readonly VITE_CLOUD_API_BASE: string
  readonly VITE_DOWNLOAD_WIN_URL: string
  readonly VITE_DOWNLOAD_VERSION: string
  readonly VITE_DOWNLOAD_USE_DIRECT: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
