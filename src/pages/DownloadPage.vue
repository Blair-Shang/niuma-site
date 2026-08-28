<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RsButton } from '@/ui'
import { fetchDownloadStats, startWindowsDownload } from '../api/downloads'
import { fetchLatestRelease, type UpdateRelease } from '../api/updates'
import { siteConfig } from '../config/site'

const total = ref<number | null>(null)
const release = ref<UpdateRelease | null>(null)
const loadingRelease = ref(true)

const versionLabel = computed(() => {
  if (release.value?.version) return `v${release.value.version}`
  if (siteConfig.download.version) return `v${siteConfig.download.version}`
  return '最新版本'
})

const notesTitle = computed(() => release.value?.title || '')
const notesBody = computed(() => (release.value?.notesMd || '').trim())

onMounted(async () => {
  const [stats, latest] = await Promise.all([
    fetchDownloadStats(),
    fetchLatestRelease({ platform: 'windows', arch: 'x64' }),
  ])
  if (stats && stats.total >= 0) total.value = stats.total
  release.value = latest
  loadingRelease.value = false
})
</script>

<template>
  <div class="site-section">
    <div class="site-container">
      <header class="site-page-head">
        <p class="site-section__eyebrow">Download</p>
        <h1>下载 NiuMa</h1>
        <p>获取桌面端安装包。通过本页下载会计入累计次数；版本与更新说明来自云端发布。</p>
        <p v-if="total !== null" class="site-stat" style="margin-top: 1.1rem">
          <span class="site-stat__value">{{ total.toLocaleString('zh-CN') }}</span>
          <span class="site-stat__label">次累计下载</span>
        </p>
      </header>

      <div class="download-grid">
        <article class="download-tile download-tile--primary">
          <div class="download-tile__head">
            <h2>Windows</h2>
            <p>x64 安装包 · {{ versionLabel }}</p>
          </div>
          <RsButton
            variant="primary"
            radius="md"
            icon="download"
            @click="startWindowsDownload"
          >
            下载 Setup.exe
          </RsButton>
          <p class="download-tile__hint">Windows 10 / 11（x64）</p>
        </article>

        <article class="download-tile">
          <div class="download-tile__head">
            <h2>Linux / macOS</h2>
            <p>将随发布流水线陆续提供。</p>
          </div>
          <span class="download-tile__soon">即将推出</span>
        </article>
      </div>

      <section class="download-notes" aria-labelledby="download-notes-title">
        <h2 id="download-notes-title">更新说明</h2>
        <p v-if="loadingRelease" class="download-notes__muted">加载中…</p>
        <template v-else-if="release">
          <p v-if="notesTitle" class="download-notes__title">{{ notesTitle }}</p>
          <pre v-if="notesBody" class="download-notes__body">{{ notesBody }}</pre>
          <p v-else class="download-notes__muted">本版本暂无详细说明。</p>
        </template>
        <p v-else class="download-notes__muted">
          暂无已发布版本说明
          <span v-if="siteConfig.download.version">
            （展示回落 {{ siteConfig.download.version }}）
          </span>
          。
        </p>
      </section>
    </div>
  </div>
</template>

<style scoped>
.download-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

@media (max-width: 800px) {
  .download-grid {
    grid-template-columns: 1fr;
  }
}

.download-tile {
  display: grid;
  gap: 1rem;
  align-content: start;
  min-height: 12rem;
  padding: 1.5rem;
  border-radius: 1rem;
  border: 1px solid color-mix(in srgb, var(--rs-border) 90%, transparent);
  background: linear-gradient(165deg, #14171f, #0e1014);
}

.download-tile--primary {
  border-color: color-mix(in srgb, var(--site-aurora-a) 28%, var(--rs-border));
  box-shadow: 0 0 0 1px color-mix(in srgb, var(--site-aurora-a) 12%, transparent) inset;
}

.download-tile__head h2 {
  margin: 0 0 0.4rem;
  font-family: var(--site-font-display);
  font-size: 1.25rem;
  letter-spacing: -0.03em;
}

.download-tile__head p,
.download-tile__hint,
.download-tile__soon {
  margin: 0;
  color: var(--rs-muted);
  line-height: 1.55;
  font-size: 0.95rem;
}

.download-notes {
  margin-top: 2.5rem;
  display: grid;
  gap: 0.75rem;
}

.download-notes h2 {
  margin: 0;
  font-family: var(--site-font-display);
  font-size: 1.35rem;
  letter-spacing: -0.03em;
}

.download-notes__title {
  margin: 0;
  font-weight: 600;
}

.download-notes__body {
  margin: 0;
  padding: 1rem 1.15rem;
  border-radius: 0.85rem;
  border: 1px solid color-mix(in srgb, var(--rs-border) 90%, transparent);
  background: color-mix(in srgb, #14171f 92%, transparent);
  white-space: pre-wrap;
  font-family: inherit;
  font-size: 0.92rem;
  line-height: 1.55;
  color: var(--rs-fg, #e8eaed);
  max-height: 22rem;
  overflow: auto;
}

.download-notes__muted {
  margin: 0;
  color: var(--rs-muted);
  font-size: 0.95rem;
}
</style>
