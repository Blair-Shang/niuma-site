<script setup lang="ts">
import { RsButton } from '@/ui'
import { useRouter } from 'vue-router'
import { siteConfig } from '../config/site'

const router = useRouter()
</script>

<template>
  <div class="site-section">
    <div class="site-container">
      <header class="site-page-head">
        <p class="site-section__eyebrow">About</p>
        <h1>关于 {{ siteConfig.name }}</h1>
        <p>
          {{ siteConfig.name }}
          是工作室品牌与旗舰产品共用名。我们专注专业开发与运维工具，官网承载品牌与产品矩阵，桌面端在独立工程中持续交付。
        </p>
      </header>

      <div class="site-prose">
        <h2>我们做什么</h2>
        <p>
          以桌面工作台形态，把数据库管理、查询、监控与 AI 辅助能力整合进统一体验，降低工具切换成本。
        </p>
        <h2>联系与反馈</h2>
        <p>缺陷与功能建议请前往反馈页；商务合作可后续在此补充邮箱。</p>
        <div class="site-cta-row">
          <RsButton variant="primary" radius="md" @click="router.push('/feedback')">
            问题反馈
          </RsButton>
          <RsButton variant="secondary" radius="md" @click="router.push('/download')">
            下载产品
          </RsButton>
        </div>
      </div>
    </div>
  </div>
</template>
