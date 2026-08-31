<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { RsBadge, RsButton } from '@/ui'
import FaqList from '../components/FaqList.vue'
import { faqItemsByCategory } from '../data/faq'
import { fetchDownloadStats, startPlatformDownload, type DownloadPlatform } from '../api/downloads'
import {
  fetchLatestRelease,
  fetchPublishedRelease,
  fetchReleaseHistory,
  type UpdateRelease,
} from '../api/updates'
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

const installFaq = faqItemsByCategory('install').slice(0, 3)
const router = useRouter()

const total = ref<number | null>(null)
const releases = ref<Partial<Record<DownloadPlatform, UpdateRelease | null>>>({})
const loadingRelease = ref(true)

const notesHistory = ref<UpdateRelease[]>([])
const selectedVersion = ref('')
const notesCache = ref<Record<string, UpdateRelease>>({})
const notesLoading = ref(false)
const copiedSha = ref('')

const selectedCard = computed(() => {
  return (
    notesHistory.value.find((rel) => rel.version === selectedVersion.value) ||
    notesHistory.value[0] ||
    null
  )
})

const selectedRelease = computed(() => {
  const ver = selectedCard.value?.version
  if (!ver) return null
  return notesCache.value[ver] || selectedCard.value
})

async function loadSelectedNotes(version: string) {
  if (!version || notesCache.value[version]?.notesMd) return
  notesLoading.value = true
  const full = await fetchPublishedRelease({
    platform: 'windows',
    arch: 'x64',
    channel: 'stable',
    version,
  })
  if (full) {
    notesCache.value = { ...notesCache.value, [version]: full }
  }
  notesLoading.value = false
}

watch(selectedVersion, (version) => {
  if (version) void loadSelectedNotes(version)
})

function formatReleaseDate(iso?: string): string {
  if (!iso) return ''
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

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
    fetchReleaseHistory({ platform: 'windows', arch: 'x64', channel: 'stable', limit: 50 }),
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
  selectedVersion.value = notesHistory.value[0]?.version || ''
  if (next.windows?.notesMd && next.windows.version) {
    notesCache.value = { ...notesCache.value, [next.windows.version]: next.windows }
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

      <section class="download-faq" aria-labelledby="download-faq-title">
        <div class="download-faq__head">
          <div>
            <h2 id="download-faq-title">安装遇到问题？</h2>
            <p class="download-faq__lead">
              当前 Windows 安装包尚未购买企业代码签名证书，SmartScreen 可能提示「未知发布者」。从本页下载并对照 SHA-256 校验后可继续安装。
            </p>
          </div>
          <RsButton variant="secondary" radius="md" @click="router.push('/faq?category=install')">
            查看全部帮助
          </RsButton>
        </div>
        <FaqList :items="installFaq" open-id="windows-smartscreen" />
      </section>

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
        <h2 id="download-notes-title">更新日志</h2>
        <p v-if="loadingRelease" class="download-notes__muted">加载中…</p>
        <div v-else-if="notesHistory.length" class="download-notes__split">
          <nav class="download-notes__nav" aria-label="版本列表">
            <button
              v-for="rel in notesHistory"
              :key="rel.version"
              type="button"
              class="download-notes__item"
              :class="{ 'is-active': selectedCard?.version === rel.version }"
              @click="selectedVersion = rel.version"
            >
              <span class="download-notes__ver-id">v{{ rel.version }}</span>
              <span class="download-notes__time">{{ formatReleaseDate(rel.publishedAt) || '—' }}</span>
            </button>
          </nav>
          <article v-if="selectedRelease" class="download-notes__detail">
            <p class="download-notes__title">
              <span class="download-notes__ver-id">v{{ selectedRelease.version }}</span>
              <template v-if="selectedRelease.title"> {{ selectedRelease.title }}</template>
            </p>
            <p v-if="notesLoading && !selectedRelease.notesMd?.trim()" class="download-notes__muted">
              加载说明中…
            </p>
            <pre v-else-if="selectedRelease.notesMd?.trim()" class="download-notes__body">{{
              selectedRelease.notesMd.trim()
            }}</pre>
            <p v-else class="download-notes__muted">本版本暂无详细说明。</p>
          </article>
        </div>
        <p v-else class="download-notes__muted">
          暂无更新说明。
        </p>
      </section>
    </div>
  </div>
</template>

<style scoped>
.download-faq {
  margin-bottom: 2.5rem;
  display: grid;
  gap: 1rem;
}

.download-faq__head {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
}

.download-faq h2 {
  margin: 0 0 0.45rem;
  font-family: var(--site-font-display);
  font-size: 1.35rem;
  letter-spacing: -0.03em;
}

.download-faq__lead {
  margin: 0;
  max-width: 44rem;
  color: var(--rs-muted);
  line-height: 1.65;
  font-size: 0.95rem;
}

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

.download-notes__split {
  display: grid;
  grid-template-columns: 13.5rem minmax(0, 1fr);
  min-height: 22rem;
  border-radius: 0.85rem;
  border: 1px solid color-mix(in srgb, var(--rs-border) 90%, transparent);
  background: color-mix(in srgb, #14171f 92%, transparent);
  overflow: hidden;
}

.download-notes__nav {
  display: flex;
  flex-direction: column;
  overflow: auto;
  border-right: 1px solid color-mix(in srgb, var(--rs-border) 90%, transparent);
  max-height: 28rem;
}

.download-notes__item {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.2rem;
  width: 100%;
  margin: 0;
  padding: 0.8rem 1rem;
  border: none;
  border-bottom: 1px solid color-mix(in srgb, var(--rs-border) 55%, transparent);
  background: transparent;
  color: inherit;
  font: inherit;
  text-align: left;
  cursor: pointer;
}

.download-notes__item:hover {
  background: color-mix(in srgb, var(--site-aurora-a, #7aa2ff) 8%, transparent);
}

.download-notes__item.is-active {
  background: color-mix(in srgb, var(--site-aurora-a, #7aa2ff) 16%, transparent);
}

.download-notes__time {
  color: var(--rs-muted);
  font-size: 0.8rem;
}

.download-notes__detail {
  display: grid;
  align-content: start;
  gap: 0.7rem;
  min-width: 0;
  padding: 1.1rem 1.2rem;
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
  max-height: 22rem;
  overflow: auto;
  white-space: pre-wrap;
  font-family: inherit;
  font-size: 0.92rem;
  line-height: 1.55;
  color: var(--rs-fg, #e8eaed);
}

.download-notes__muted {
  margin: 0;
  color: var(--rs-muted);
  font-size: 0.95rem;
}

@media (max-width: 800px) {
  .download-notes__split {
    grid-template-columns: 1fr;
  }

  .download-notes__nav {
    flex-direction: row;
    max-height: none;
    overflow: auto;
    border-right: none;
    border-bottom: 1px solid color-mix(in srgb, var(--rs-border) 90%, transparent);
  }

  .download-notes__item {
    flex: 0 0 auto;
    min-width: 7.5rem;
    border-bottom: none;
    border-right: 1px solid color-mix(in srgb, var(--rs-border) 55%, transparent);
  }
}
</style>
