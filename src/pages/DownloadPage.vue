<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RsBadge, RsButton } from '@/ui'
import { fetchDownloadStats, startPlatformDownload, type DownloadPlatform } from '../api/downloads'
import { fetchLatestRelease, fetchReleaseHistory, type UpdateRelease } from '../api/updates'
import { siteConfig } from '../config/site'

type PlatformCard = {
  id: DownloadPlatform
  title: string
  arch: string
  channel: 'stable' | 'beta'
  button: string
  hint: string
}

const cards: PlatformCard[] = [
  {
    id: 'windows',
    title: 'Windows',
    arch: 'x64',
    channel: 'stable',
    button: '下载 Setup.exe',
    hint: 'Windows 10 / 11（x64）',
  },
  {
    id: 'linux',
    title: 'Linux',
    arch: 'x64',
    channel: 'beta',
    button: '下载预览版',
    hint: 'x64 · .run / .deb，预览版',
  },
  {
    id: 'macos',
    title: 'macOS',
    arch: 'arm64',
    channel: 'beta',
    button: '下载预览版',
    hint: 'Apple Silicon · .pkg / .dmg，预览版',
  },
]

const total = ref<number | null>(null)
const releases = ref<Partial<Record<DownloadPlatform, UpdateRelease | null>>>({})
const loadingRelease = ref(true)

const notesHistory = ref<UpdateRelease[]>([])
const copiedSha = ref('')

function sha256Of(id: DownloadPlatform): string {
  return (releases.value[id]?.sha256 || '').trim().toLowerCase()
}

function shaShort(hex: string): string {
  if (hex.length < 16) return ''
  return `${hex.slice(0, 12)}…${hex.slice(-8)}`
}

async function copySha(id: DownloadPlatform) {
  const hex = sha256Of(id)
  if (!hex) return
  try {
    await navigator.clipboard.writeText(hex)
    copiedSha.value = id
    window.setTimeout(() => {
      if (copiedSha.value === id) copiedSha.value = ''
    }, 1600)
  } catch {
    /* ignore */
  }
}

function versionLabel(id: DownloadPlatform): string {
  const rel = releases.value[id]
  if (rel?.version) return `v${rel.version}`
  if (id === 'windows' && siteConfig.download.version) return `v${siteConfig.download.version}`
  return '暂无发布'
}

function hasPackage(id: DownloadPlatform): boolean {
  return !!releases.value[id]?.version
}

onMounted(async () => {
  const statsP = fetchDownloadStats()
  const rows = await Promise.all(
    cards.map(async (c) => [c.id, await fetchLatestRelease({ platform: c.id, arch: c.arch, channel: c.channel })] as const),
  )
  const next: Partial<Record<DownloadPlatform, UpdateRelease | null>> = {}
  for (const [id, rel] of rows) next[id] = rel
  const [stats, history] = await Promise.all([
    statsP,
    fetchReleaseHistory({ platform: 'windows', arch: 'x64', channel: 'stable', limit: 10 }),
  ])
  if (stats && stats.total >= 0) total.value = stats.total
  releases.value = next
  if (history.length > 0) {
    notesHistory.value = history
  } else if (next.windows) {
    notesHistory.value = [next.windows]
  } else {
    notesHistory.value = []
  }
  loadingRelease.value = false
})
</script>

<template>
  <div class="site-section">
    <div class="site-container">
      <header class="site-page-head">
        <p class="site-section__eyebrow">Download</p>
        <h1>下载 NiuMa</h1>
        <p>
          推荐下载 Windows 正式版。Linux 与 macOS 目前为预览版，功能与稳定性仍在完善中。
        </p>
        <p v-if="total !== null" class="site-stat" style="margin-top: 1.1rem">
          <span class="site-stat__value">{{ total.toLocaleString('zh-CN') }}</span>
          <span class="site-stat__label">次累计下载（含应用内更新）</span>
        </p>
      </header>

      <div class="download-grid">
        <article
          v-for="card in cards"
          :key="card.id"
          class="download-tile"
          :class="{ 'download-tile--primary': card.channel === 'stable' }"
        >
          <div class="download-tile__head">
            <h2>
              {{ card.title }}
              <RsBadge v-if="card.channel === 'beta'" variant="warning">Beta</RsBadge>
            </h2>
            <p>{{ card.arch }} 安装包 · {{ versionLabel(card.id) }}</p>
          </div>
          <RsButton
            :variant="card.channel === 'stable' ? 'primary' : 'secondary'"
            radius="md"
            icon="download"
            :disabled="!hasPackage(card.id)"
            @click="startPlatformDownload(card.id)"
          >
            {{ hasPackage(card.id) ? card.button : card.channel === 'beta' ? '暂无预览包' : '暂无安装包' }}
          </RsButton>
          <p class="download-tile__hint">{{ card.hint }}</p>
          <p v-if="sha256Of(card.id)" class="download-tile__sha">
            <span class="download-tile__sha-label">SHA-256</span>
            <button type="button" class="download-tile__sha-btn" :title="sha256Of(card.id)" @click="copySha(card.id)">
              {{ copiedSha === card.id ? '已复制' : shaShort(sha256Of(card.id)) }}
            </button>
          </p>
        </article>
      </div>

      <section class="download-notes" aria-labelledby="download-notes-title">
        <h2 id="download-notes-title">更新说明</h2>
        <p v-if="loadingRelease" class="download-notes__muted">加载中…</p>
        <template v-else-if="notesHistory.length">
          <article v-for="rel in notesHistory" :key="rel.version" class="download-notes__ver">
            <p class="download-notes__title">
              <span class="download-notes__ver-id">v{{ rel.version }}</span>
              <template v-if="rel.title"> {{ rel.title }}</template>
            </p>
            <pre v-if="rel.notesMd?.trim()" class="download-notes__body">{{ rel.notesMd.trim() }}</pre>
            <p v-else class="download-notes__muted">本版本暂无详细说明。</p>
          </article>
        </template>
        <p v-else class="download-notes__muted">
          暂无更新说明。
        </p>
      </section>
    </div>
  </div>
</template>

<style scoped>
.download-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1rem;
}

@media (max-width: 960px) {
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
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-family: var(--site-font-display);
  font-size: 1.25rem;
  letter-spacing: -0.03em;
}

.download-tile__head p,
.download-tile__hint {
  margin: 0;
  color: var(--rs-muted);
  line-height: 1.55;
  font-size: 0.95rem;
}

.download-tile__sha {
  margin: 0;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.45rem;
  font-size: 0.78rem;
  color: var(--rs-muted);
}

.download-tile__sha-label {
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.download-tile__sha-btn {
  margin: 0;
  padding: 0;
  border: none;
  background: none;
  color: inherit;
  font: inherit;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  cursor: pointer;
  text-decoration: underline;
  text-underline-offset: 0.15em;
}

.download-tile__sha-btn:hover {
  color: var(--rs-fg, #e8eaed);
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

.download-notes__ver {
  display: grid;
  gap: 0.5rem;
}

.download-notes__title {
  margin: 0;
  font-weight: 600;
}

.download-notes__ver-id {
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-weight: 650;
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
