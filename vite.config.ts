import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { resolve } from 'node:path'

const niumaUiRoot = resolve(__dirname, '../niuma-ui')
const niumaUiSrc = resolve(niumaUiRoot, 'src')

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: [
      {
        find: '@niuma/ui/styles.css',
        replacement: resolve(niumaUiSrc, 'styles.css'),
      },
      // 组件按需从源码目录引用（见 src/ui.ts），不要指向整包 index.ts
      {
        find: '@niuma-ui-src',
        replacement: niumaUiSrc,
      },
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
})
