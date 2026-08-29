<script setup lang="ts">
import { useRouter } from 'vue-router'
import { siteConfig } from '../config/site'

const router = useRouter()
const year = new Date().getFullYear()

const columns = [
  {
    title: '产品',
    links: [
      { label: '产品矩阵', to: '/products' },
      { label: 'NiuMa 桌面端', to: '/products/niuma' },
      { label: '下载', to: '/download' },
    ],
  },
  {
    title: '支持',
    links: [
      { label: '问题反馈', to: '/feedback' },
      { label: '关于', to: '/about' },
    ],
  },
] as const
</script>

<template>
  <footer class="site-footer">
    <div class="site-container site-footer__grid">
      <div class="site-footer__brand-col">
        <div class="site-footer__brand-row">
          <img src="/brand/app-icon.svg" alt="" width="28" height="28" />
          <span>{{ siteConfig.name }}</span>
        </div>
        <p class="site-footer__tag">{{ siteConfig.tagline }}</p>
      </div>

      <div
        v-for="col in columns"
        :key="col.title"
        class="site-footer__col"
      >
        <p class="site-footer__col-title">{{ col.title }}</p>
        <button
          v-for="link in col.links"
          :key="link.to"
          type="button"
          class="site-footer__link"
          @click="router.push(link.to)"
        >
          {{ link.label }}
        </button>
      </div>
    </div>
    <div class="site-container site-footer__bottom">
      <p class="site-footer__copy">© {{ year }} {{ siteConfig.name }}</p>
      <a
        class="site-footer__beian"
        :href="siteConfig.icp.url"
        target="_blank"
        rel="noopener noreferrer"
      >
        {{ siteConfig.icp.number }}
      </a>
    </div>
  </footer>
</template>

<style scoped>
.site-footer {
  margin-top: auto;
  border-top: 1px solid color-mix(in srgb, var(--rs-border) 85%, transparent);
  padding: 3rem 0 2rem;
  background: color-mix(in srgb, #08090c 90%, transparent);
}

.site-footer__grid {
  display: grid;
  gap: 2rem;
  padding-bottom: 2rem;
}

@media (min-width: 800px) {
  .site-footer__grid {
    grid-template-columns: 1.4fr 1fr 1fr;
  }
}

.site-footer__brand-row {
  display: inline-flex;
  align-items: center;
  gap: 0.55rem;
  font-family: var(--site-font-display);
  font-weight: 650;
  letter-spacing: -0.02em;
}

.site-footer__brand-row img {
  border-radius: 22%;
}

.site-footer__tag {
  margin: 0.75rem 0 0;
  max-width: 18rem;
  color: var(--rs-muted);
  line-height: 1.6;
  font-size: 0.92rem;
}

.site-footer__col {
  display: grid;
  gap: 0.55rem;
  align-content: start;
}

.site-footer__col-title {
  margin: 0 0 0.35rem;
  font-size: 0.8rem;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: color-mix(in srgb, #fff 55%, var(--rs-muted));
}

.site-footer__link {
  appearance: none;
  border: 0;
  background: transparent;
  padding: 0;
  width: fit-content;
  font: inherit;
  font-size: 0.95rem;
  color: var(--rs-muted);
  cursor: pointer;
  text-align: left;
}

.site-footer__link:hover {
  color: #fff;
}

.site-footer__bottom {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.75rem 1.25rem;
  border-top: 1px solid color-mix(in srgb, var(--rs-border) 70%, transparent);
  padding-top: 1.25rem;
}

.site-footer__copy,
.site-footer__beian {
  margin: 0;
  font-size: 0.85rem;
  color: var(--rs-muted);
}

.site-footer__beian {
  text-decoration: none;
}

.site-footer__beian:hover {
  color: #fff;
}
</style>
