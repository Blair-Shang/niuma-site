<script setup lang="ts">
import { RsButton } from '@/ui'
import { useRouter } from 'vue-router'
import ProductShotGallery from '../components/ProductShotGallery.vue'
import { niumaDesktopShots } from '../data/products'
import { siteConfig } from '../config/site'

const router = useRouter()

const highlights = [
  {
    title: '数据库工作台',
    description: 'MySQL、Vastbase 等会话、对象树、查询与表设计在同一工作台完成。',
  },
  {
    title: 'AI 辅助运维',
    description: 'Skills 与工具编排结合会话策略，缩短排障与变更路径。',
  },
  {
    title: '桌面原生体验',
    description: '本地安装、核心能力可离线使用，适配内网与合规场景。',
  },
] as const
</script>

<template>
  <div class="site-section">
    <div class="site-container">
      <header class="site-page-head">
        <p class="site-section__eyebrow">Product</p>
        <h1>NiuMa 桌面端</h1>
        <p>
          全能 AI 运维平台：连接管理、SQL、数据浏览、监控与外部工具链收敛到同一桌面工作台。
        </p>
        <div class="site-cta-row">
          <RsButton
            variant="primary"
            radius="md"
            icon="download"
            @click="router.push('/download')"
          >
            前往下载
          </RsButton>
          <RsButton
            variant="secondary"
            radius="md"
            @click="router.push('/feedback')"
          >
            问题反馈
          </RsButton>
        </div>
      </header>

      <ProductShotGallery
        class="product-shot"
        :shots="niumaDesktopShots"
      />

      <div class="site-pillars">
        <article
          v-for="item in highlights"
          :key="item.title"
          class="site-pillar"
        >
          <h3 class="site-pillar__title">{{ item.title }}</h3>
          <p class="site-pillar__desc">{{ item.description }}</p>
        </article>
      </div>

      <p class="product-foot">
        {{ siteConfig.name }} 持续迭代中。欢迎通过反馈渠道提交缺陷与需求。
      </p>
    </div>
  </div>
</template>

<style scoped>
.product-shot {
  margin: 0 0 2rem;
}

.product-foot {
  margin: 2rem 0 0;
  font-size: 0.9375rem;
  color: var(--rs-muted);
}
</style>
