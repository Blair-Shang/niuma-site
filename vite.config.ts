import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { existsSync } from 'node:fs'
import { createRequire } from 'node:module'
import { dirname, resolve } from 'node:path'

const require = createRequire(import.meta.url)

/** 本机有同级源码则联调；CI / 打包从 npm 的 niuma-ui 解析。 */
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
const niumaUiSrc = resolve(niumaUiRoot, 'src')
const niumaUiDist = resolve(niumaUiRoot, 'dist/index.js')

/** 本机 dev 走兄弟仓源码 HMR；生产构建优先 dist。 */
function niumaUiAliases(command: string) {
  if (!existsSync(niumaUiSrc)) return []
  if (command !== 'serve' && existsSync(niumaUiDist)) return []
  return [
    {
      find: '@niuma/ui/styles.css',
      replacement: resolve(niumaUiSrc, 'styles.css'),
    },
    {
      find: /^@niuma\/ui$/,
      replacement: resolve(niumaUiSrc, 'index.ts'),
    },
  ]
}

export default defineConfig(({ command }) => ({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: [
      ...niumaUiAliases(command),
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
    cssCodeSplit: true,
  },
}))
