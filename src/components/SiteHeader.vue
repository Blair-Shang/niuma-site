<script setup lang="ts">
import { RsButton } from '@/ui'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()

const nav = [
  { label: '产品', to: '/#products' },
  { label: '下载', to: '/download' },
  { label: '反馈', to: '/feedback' },
  { label: '关于', to: '/about' },
] as const

function go(path: string) {
  if (path.startsWith('/#')) {
    void router.push({ path: '/', hash: path.slice(1) })
    return
  }
  void router.push(path)
}

function isActive(to: string) {
  if (to.startsWith('/#')) return route.path === '/'
  return route.path === to || route.path.startsWith(`${to}/`)
}
</script>

<template>
  <header class="site-header">
    <div class="site-container site-header__shell">
      <RouterLink to="/" class="site-header__brand">
        <img src="/brand/app-icon.svg" alt="" width="26" height="26" />
        <span>NiuMa</span>
      </RouterLink>

      <nav class="site-header__nav" aria-label="主导航">
        <button
          v-for="item in nav"
          :key="item.to"
          type="button"
          class="site-header__link"
          :class="{ 'site-header__link--active': isActive(item.to) }"
          @click="go(item.to)"
        >
          {{ item.label }}
        </button>
      </nav>

      <div class="site-header__cta">
        <RsButton
          variant="primary"
          size="sm"
          radius="md"
          icon="download"
          @click="go('/download')"
        >
          下载
        </RsButton>
      </div>
    </div>
  </header>
</template>

<style scoped>
.site-header {
  position: sticky;
  top: 0;
  z-index: 40;
  padding: 0.85rem 0 0.65rem;
}

.site-header__shell {
  display: flex;
  align-items: center;
  gap: 1rem;
  min-height: 3.25rem;
  padding: 0.35rem 0.55rem 0.35rem 0.85rem;
  border: 1px solid color-mix(in srgb, var(--rs-border) 80%, transparent);
  border-radius: 999px;
  background: color-mix(in srgb, #0b0c0f 72%, transparent);
  backdrop-filter: blur(18px) saturate(140%);
  -webkit-backdrop-filter: blur(18px) saturate(140%);
  box-shadow: 0 10px 40px color-mix(in srgb, #000 35%, transparent);
}

.site-header__brand {
  display: inline-flex;
  align-items: center;
  gap: 0.55rem;
  margin-right: auto;
  color: #fff;
  text-decoration: none;
  font-family: var(--site-font-display);
  font-weight: 650;
  letter-spacing: -0.03em;
  font-size: 0.98rem;
}

.site-header__brand img {
  border-radius: 22%;
}

.site-header__nav {
  display: none;
  align-items: center;
  gap: 0.1rem;
}

.site-header__link {
  appearance: none;
  border: 0;
  background: transparent;
  color: var(--rs-muted);
  font: inherit;
  font-size: 0.9rem;
  font-weight: 500;
  padding: 0.45rem 0.8rem;
  border-radius: 999px;
  cursor: pointer;
  transition: color 0.15s ease, background 0.15s ease;
}

.site-header__link:hover {
  color: #fff;
  background: color-mix(in srgb, #fff 6%, transparent);
}

.site-header__link--active {
  color: #fff;
}

.site-header__cta {
  display: flex;
  align-items: center;
}

@media (min-width: 860px) {
  .site-header__nav {
    display: flex;
  }
}
</style>
