<script setup lang="ts">
import ProductCard from '../components/ProductCard.vue'
import { products } from '../data/products'
</script>

<template>
  <div class="site-section">
    <div class="site-container">
      <header class="site-page-head">
        <p class="site-section__eyebrow">Products</p>
        <h1>产品矩阵</h1>
        <p>
          工作室能力以产品卡片扩展。当前主推桌面端，后续工具将并列进入同一入口。
        </p>
      </header>
      <div class="site-product-grid">
        <ProductCard
          v-for="item in products"
          :key="item.id"
          :product="item"
        />
      </div>
    </div>
  </div>
</template>
