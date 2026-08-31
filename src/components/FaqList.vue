<script setup lang="ts">
import { useRouter } from 'vue-router'
import type { FaqItem } from '../data/faq'

defineProps<{
  items: FaqItem[]
  /** 无匹配结果时的提示 */
  emptyText?: string
  /** 默认展开的首条 id */
  openId?: string
}>()

const router = useRouter()

function onBodyClick(event: MouseEvent) {
  const anchor = (event.target as HTMLElement | null)?.closest('a')
  if (!anchor) return
  const href = anchor.getAttribute('href')
  if (!href || !href.startsWith('/') || href.startsWith('//')) return
  event.preventDefault()
  void router.push(href)
}
</script>

<template>
  <div v-if="items.length" class="faq-list">
    <details
      v-for="item in items"
      :key="item.id"
      :id="item.id"
      class="faq-list__item"
      :open="openId === item.id"
    >
      <summary>{{ item.question }}</summary>
      <div class="faq-list__body" @click="onBodyClick">
        <template v-for="(block, idx) in item.answer" :key="`${item.id}-${idx}`">
          <p v-if="block.kind === 'text'" class="faq-list__text" v-html="block.html" />
          <ol v-else-if="block.kind === 'list' && block.ordered" class="faq-list__list">
            <li v-for="(line, i) in block.items" :key="i" v-html="line" />
          </ol>
          <ul v-else class="faq-list__list">
            <li v-for="(line, i) in block.items" :key="i" v-html="line" />
          </ul>
        </template>
      </div>
    </details>
  </div>
  <p v-else class="faq-list__empty">{{ emptyText ?? '没有匹配的问题。' }}</p>
</template>

<style scoped>
.faq-list {
  display: grid;
  gap: 0.55rem;
}

.faq-list__item {
  border-radius: 0.75rem;
  border: 1px solid color-mix(in srgb, var(--rs-border) 90%, transparent);
  background: color-mix(in srgb, #14171f 88%, transparent);
  overflow: hidden;
  scroll-margin-top: 5.5rem;
}

.faq-list__item summary {
  padding: 0.85rem 1rem;
  font-weight: 600;
  cursor: pointer;
  list-style: none;
}

.faq-list__item summary::-webkit-details-marker {
  display: none;
}

.faq-list__item summary::after {
  content: '+';
  float: right;
  color: var(--rs-muted);
  font-weight: 400;
}

.faq-list__item[open] summary::after {
  content: '−';
}

.faq-list__body {
  padding: 0 1rem 1rem;
  color: var(--rs-muted);
  font-size: 0.92rem;
  line-height: 1.65;
}

.faq-list__text {
  margin: 0 0 0.65rem;
}

.faq-list__text:last-child {
  margin-bottom: 0;
}

.faq-list__list {
  margin: 0 0 0.65rem;
  padding-left: 1.25rem;
}

.faq-list__list:last-child {
  margin-bottom: 0;
}

.faq-list__body :deep(code) {
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 0.85em;
}

.faq-list__body :deep(a) {
  color: color-mix(in srgb, var(--site-aurora-a) 80%, #fff);
  text-decoration: underline;
  text-underline-offset: 0.15em;
}

.faq-list__empty {
  margin: 0;
  color: var(--rs-muted);
  font-size: 0.95rem;
}
</style>
