<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RsButton } from '@/ui'
import { useRouter } from 'vue-router'
import ProductCard from '../components/ProductCard.vue'
import { fetchDownloadStats } from '../api/downloads'
import { products } from '../data/products'
import { siteConfig } from '../config/site'

const router = useRouter()
const downloadTotal = ref<number | null>(null)

const pillars = [
  {
    title: '统一工作台',
    description: '连接、对象树、查询与设计在同一桌面空间完成，减少工具切换。',
  },
  {
    title: 'AI 可编排',
    description: 'Skills 与工具链结合会话策略，让排障与变更路径更短。',
  },
  {
    title: '本地可信',
    description: '桌面安装、核心能力可离线，适配内网与合规场景。',
  },
] as const

onMounted(async () => {
  const stats = await fetchDownloadStats()
  if (stats && stats.total > 0) downloadTotal.value = stats.total
})
</script>

<template>
  <div>
    <section class="site-hero" aria-labelledby="site-hero-title">
      <div class="site-container">
        <img
          class="site-hero__mark"
          src="/brand/app-icon.svg"
          alt=""
          width="68"
          height="68"
        />
        <p class="site-hero__kicker">{{ siteConfig.tagline }}</p>
        <h1 id="site-hero-title" class="site-hero__title">NiuMa</h1>
        <p class="site-hero__lead">
          面向数据库与运维场景的 AI 桌面平台。把会话、SQL、监控与工具链收敛到一个专业工作台。
        </p>
        <div class="site-hero__actions">
          <RsButton
            variant="primary"
            size="lg"
            radius="md"
            icon="download"
            @click="router.push('/download')"
          >
            下载桌面端
          </RsButton>
          <RsButton
            variant="secondary"
            size="lg"
            radius="md"
            @click="router.push('/products/niuma')"
          >
            查看产品
          </RsButton>
        </div>
        <p v-if="downloadTotal !== null" class="site-hero__meta site-stat">
          <span class="site-stat__value">{{ downloadTotal.toLocaleString('zh-CN') }}</span>
          <span class="site-stat__label">次累计下载</span>
        </p>

        <div class="site-hero__stage" aria-hidden="true">
          <div class="site-hero__stage-grid" />
          <div class="site-hero__stage-glow" />
          <div class="site-hero__stage-card">
            <div class="site-hero__stage-bar">
              <span class="site-hero__stage-dot" />
              <span class="site-hero__stage-dot" />
              <span class="site-hero__stage-dot" />
            </div>
            <div class="site-hero__stage-body">
              <div class="site-hero__stage-side">
                <span class="site-hero__stage-line site-hero__stage-line--accent" />
                <span class="site-hero__stage-line" />
                <span class="site-hero__stage-line" />
                <span class="site-hero__stage-line" />
                <span class="site-hero__stage-line" />
              </div>
              <div class="site-hero__stage-main">
                <span class="site-hero__stage-row" />
                <span class="site-hero__stage-row" />
                <span class="site-hero__stage-row" />
                <span class="site-hero__stage-row" />
                <span class="site-hero__stage-row" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="site-section site-section--muted">
      <div class="site-container">
        <div class="site-section__head">
          <p class="site-section__eyebrow">Why NiuMa</p>
          <h2 class="site-section__title">为专业运维工作流而设计</h2>
          <p class="site-section__desc">
            不是功能堆叠的面板集合，而是一条从连接到执行、从排查到变更的连续路径。
          </p>
        </div>
        <div class="site-pillars">
          <article v-for="item in pillars" :key="item.title" class="site-pillar">
            <h3 class="site-pillar__title">{{ item.title }}</h3>
            <p class="site-pillar__desc">{{ item.description }}</p>
          </article>
        </div>
      </div>
    </section>

    <section class="site-section">
      <div class="site-container">
        <div class="site-section__head">
          <p class="site-section__eyebrow">Products</p>
          <h2 class="site-section__title">产品矩阵</h2>
          <p class="site-section__desc">
            工作室能力以产品卡片扩展。当前主推桌面端，后续工具将并列进入同一入口。
          </p>
        </div>
        <div class="site-product-grid">
          <ProductCard
            v-for="item in products"
            :key="item.id"
            :product="item"
          />
        </div>
      </div>
    </section>
  </div>
</template>
