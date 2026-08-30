import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { niumaUiHost } from '@niuma/ui/vite-plugins/niuma-ui-host'
import { existsSync, writeFileSync } from 'node:fs'
import { createRequire } from 'node:module'
import { dirname, resolve } from 'node:path'

const require = createRequire(import.meta.url)

/** 定位 @niuma/ui 包根，仅给 dev server 放行读取（link 到兄弟仓时需要）。 */
function resolveNiumaUiRoot(): string {
  const sibling = resolve(__dirname, '../niuma-ui')
  if (existsSync(resolve(sibling, 'package.json'))) {
    return sibling
  }
  try {
    return dirname(require.resolve('@niuma/ui/package.json'))
  } catch {
    return dirname(require.resolve('niuma-ui/package.json'))
  }
}

const niumaUiRoot = resolveNiumaUiRoot()

/**
 * niumaUiHost：只服务 pnpm dev（源码 HMR）。vite build 走包主入口。
 */
export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    ...niumaUiHost(),
    {
      name: 'keep-dist-gitkeep',
      closeBundle() {
        writeFileSync(resolve(__dirname, 'server/internal/web/dist/.gitkeep'), '')
      },
    },
  ],
  resolve: {
    alias: [
      {
        find: '@',
        replacement: resolve(__dirname, 'src'),
      },
    ],
  },
  server: {
    port: 5173,
    fs: {
      allow: [resolve(__dirname), niumaUiRoot],
    },
    proxy: {
      // 本地：前端 5173 → Go API 8080（需先 pnpm dev:server）
      '/api': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
      },
      '/healthz': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'server/internal/web/dist',
    emptyOutDir: true,
    target: 'es2022',
    cssCodeSplit: false,
  },
})
