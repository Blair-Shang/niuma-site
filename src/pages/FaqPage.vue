<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { RsButton, RsInput } from '@/ui'
import FaqList from '../components/FaqList.vue'
import {
  faqCategories,
  faqItems,
  faqItemsByCategory,
  faqSearchText,
  type FaqCategory,
} from '../data/faq'

const route = useRoute()
const router = useRouter()

const query = ref('')
const category = ref<FaqCategory | 'all'>('all')

const filteredItems = computed(() => {
  const base = faqItemsByCategory(category.value)
  const q = query.value.trim().toLowerCase()
  if (!q) return base
  return base.filter((item) => faqSearchText(item).includes(q))
})

const resultCount = computed(() => filteredItems.value.length)

const openId = computed(() => {
  const hash = route.hash.replace(/^#/, '')
  if (hash && filteredItems.value.some((item) => item.id === hash)) return hash
  if (query.value.trim() && filteredItems.value.length === 1) return filteredItems.value[0]?.id
  return undefined
})

function syncFromRoute() {
  const q = route.query.q
  if (typeof q === 'string') query.value = q

  const cat = route.query.category
  if (cat === 'all' || !cat) {
    category.value = 'all'
  } else if (faqCategories.some((c) => c.id === cat)) {
    category.value = cat as FaqCategory
  }
}

function pushRoute() {
  const nextQuery: Record<string, string> = {}
  const q = query.value.trim()
  if (q) nextQuery.q = q
  if (category.value !== 'all') nextQuery.category = category.value
  void router.replace({ path: '/faq', query: nextQuery, hash: route.hash || undefined })
}

watch([query, category], () => {
  pushRoute()
})

onMounted(() => {
  syncFromRoute()
  if (route.hash) {
    window.setTimeout(() => {
      document.getElementById(route.hash.slice(1))?.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }, 80)
  }
})

watch(
  () => route.fullPath,
  () => syncFromRoute(),
)
</script>

<template>
  <div class="site-section">
    <div class="site-container">
      <header class="site-page-head">
        <p class="site-section__eyebrow">Help</p>
        <h1>常见问题</h1>
        <p>
          安装、更新与使用中的常见疑问集中在此，支持关键词搜索。未找到答案请
          <RouterLink to="/feedback">提交反馈</RouterLink>。
        </p>
      </header>

      <div class="faq-toolbar">
        <label class="faq-search">
          <span class="faq-search__label">搜索</span>
          <RsInput
            v-model="query"
            placeholder="例如：SmartScreen、签名、更新、反馈"
            clearable
          />
        </label>
        <nav class="faq-tabs" aria-label="问题分类">
          <button
            type="button"
            class="faq-tabs__btn"
            :class="{ 'is-active': category === 'all' }"
            @click="category = 'all'"
          >
            全部
            <span class="faq-tabs__count">{{ faqItems.length }}</span>
          </button>
          <button
            v-for="cat in faqCategories"
            :key="cat.id"
            type="button"
            class="faq-tabs__btn"
            :class="{ 'is-active': category === cat.id }"
            @click="category = cat.id"
          >
            {{ cat.label }}
            <span class="faq-tabs__count">{{ faqItemsByCategory(cat.id).length }}</span>
          </button>
        </nav>
      </div>

      <p class="faq-meta">
        <template v-if="query.trim()">
          找到 {{ resultCount }} 条与「{{ query.trim() }}」相关的问题
        </template>
        <template v-else>共 {{ resultCount }} 条</template>
      </p>

      <FaqList
        :items="filteredItems"
        :open-id="openId"
        empty-text="没有匹配的问题，可尝试其他关键词或提交反馈。"
      />

      <div class="faq-footer">
        <p>仍未解决？</p>
        <RsButton variant="primary" radius="md" @click="router.push('/feedback')">
          问题反馈
        </RsButton>
        <RsButton variant="secondary" radius="md" @click="router.push('/download')">
          前往下载
        </RsButton>
      </div>
    </div>
  </div>
</template>

<style scoped>
.faq-toolbar {
  display: grid;
  gap: 1rem;
  margin-bottom: 0.75rem;
}

.faq-search {
  display: grid;
  gap: 0.45rem;
  max-width: 28rem;
}

.faq-search__label {
  font-size: 0.9rem;
  color: var(--rs-muted);
}

.faq-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
}

.faq-tabs__btn {
  appearance: none;
  border: 1px solid color-mix(in srgb, var(--rs-border) 90%, transparent);
  background: color-mix(in srgb, #14171f 88%, transparent);
  color: var(--rs-muted);
  font: inherit;
  font-size: 0.88rem;
  padding: 0.45rem 0.75rem;
  border-radius: 999px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  transition: color 0.15s ease, border-color 0.15s ease, background 0.15s ease;
}

.faq-tabs__btn:hover {
  color: #fff;
}

.faq-tabs__btn.is-active {
  color: #fff;
  border-color: color-mix(in srgb, var(--site-aurora-a) 35%, var(--rs-border));
  background: color-mix(in srgb, var(--site-aurora-a) 14%, #14171f);
}

.faq-tabs__count {
  font-size: 0.78rem;
  color: color-mix(in srgb, var(--rs-muted) 85%, #fff);
  font-variant-numeric: tabular-nums;
}

.faq-tabs__btn.is-active .faq-tabs__count {
  color: color-mix(in srgb, #fff 70%, var(--site-aurora-a));
}

.faq-meta {
  margin: 0 0 0.85rem;
  color: var(--rs-muted);
  font-size: 0.88rem;
}

.faq-footer {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.65rem 0.75rem;
  margin-top: 2.25rem;
  padding-top: 1.5rem;
  border-top: 1px solid color-mix(in srgb, var(--rs-border) 70%, transparent);
}

.faq-footer p {
  margin: 0;
  width: 100%;
  color: var(--rs-muted);
  font-size: 0.95rem;
}

.site-page-head :deep(a) {
  color: color-mix(in srgb, var(--site-aurora-a) 80%, #fff);
  text-decoration: underline;
  text-underline-offset: 0.15em;
}
</style>
