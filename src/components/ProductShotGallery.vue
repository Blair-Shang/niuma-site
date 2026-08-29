<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { RsButton } from '@/ui'
import type { ProductShot } from '../data/products'

defineOptions({ inheritAttrs: false })

const props = defineProps<{
  shots: readonly ProductShot[]
}>()

const active = ref<number | null>(null)
const actualSize = ref(false)

const current = computed(() =>
  active.value === null ? null : (props.shots[active.value] ?? null),
)

function openAt(index: number) {
  active.value = index
  actualSize.value = false
}

function close() {
  active.value = null
  actualSize.value = false
}

function openOriginal() {
  if (!current.value) return
  window.open(current.value.src, '_blank', 'noopener,noreferrer')
}

function step(delta: number) {
  if (active.value === null || props.shots.length < 2) return
  const n = props.shots.length
  active.value = (active.value + delta + n) % n
}

function onKeydown(ev: KeyboardEvent) {
  if (active.value === null) return
  if (ev.key === 'Escape') {
    ev.preventDefault()
    close()
    return
  }
  if (ev.key === 'ArrowLeft') {
    ev.preventDefault()
    step(-1)
    return
  }
  if (ev.key === 'ArrowRight') {
    ev.preventDefault()
    step(1)
  }
}

watch(active, (index, prev) => {
  const open = index !== null
  if (open === (prev !== null)) return
  document.documentElement.style.overflow = open ? 'hidden' : ''
  if (open) window.addEventListener('keydown', onKeydown)
  else window.removeEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  document.documentElement.style.overflow = ''
  window.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <div class="site-shot-gallery" v-bind="$attrs">
    <figure
      v-for="(shot, index) in shots"
      :key="shot.src"
      class="site-shot"
    >
      <button
        type="button"
        class="site-shot__trigger"
        :aria-label="`查看原图：${shot.caption}`"
        @click="openAt(index)"
      >
        <img
          class="site-shot__img"
          :src="shot.src"
          :alt="shot.alt"
          :width="shot.width"
          :height="shot.height"
        />
        <span class="site-shot__hint">点击查看原图</span>
      </button>
      <figcaption class="site-shot__caption">{{ shot.caption }}</figcaption>
    </figure>
  </div>

  <Teleport to="body">
    <div
      v-if="current"
      class="site-lightbox"
      role="dialog"
      aria-modal="true"
      :aria-label="current.caption"
    >
      <button
        type="button"
        class="site-lightbox__backdrop"
        tabindex="-1"
        aria-label="关闭预览"
        @click="close"
      />
      <div class="site-lightbox__chrome">
        <p class="site-lightbox__title">{{ current.caption }}</p>
        <div class="site-lightbox__actions">
          <RsButton
            variant="secondary"
            size="sm"
            radius="md"
            :aria-pressed="actualSize"
            @click="actualSize = !actualSize"
          >
            {{ actualSize ? '适应窗口' : '1:1 原图' }}
          </RsButton>
          <RsButton
            variant="secondary"
            size="sm"
            radius="md"
            icon="external-link"
            @click="openOriginal"
          >
            打开原图
          </RsButton>
          <RsButton
            variant="secondary"
            size="sm"
            radius="md"
            icon="x"
            icon-only
            aria-label="关闭预览"
            @click="close"
          />
        </div>
      </div>
      <div
        class="site-lightbox__stage"
        :class="{ 'site-lightbox__stage--actual': actualSize }"
      >
        <button
          v-if="shots.length > 1"
          type="button"
          class="site-lightbox__nav site-lightbox__nav--prev"
          aria-label="上一张"
          @click="step(-1)"
        >
          ‹
        </button>
        <img
          class="site-lightbox__img"
          :class="{ 'site-lightbox__img--actual': actualSize }"
          :src="current.src"
          :alt="current.alt"
          :width="current.width"
          :height="current.height"
          :style="{
            '--shot-w': `${current.width}px`,
            '--shot-h': `${current.height}px`,
          }"
        />
        <button
          v-if="shots.length > 1"
          type="button"
          class="site-lightbox__nav site-lightbox__nav--next"
          aria-label="下一张"
          @click="step(1)"
        >
          ›
        </button>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.site-shot__trigger {
  position: relative;
  display: block;
  width: 100%;
  padding: 0;
  border: 0;
  background: #0b0c0f;
  cursor: zoom-in;
  text-align: left;
}

.site-shot__trigger:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--site-aurora-a) 70%, transparent);
  outline-offset: -2px;
}

.site-shot__hint {
  position: absolute;
  right: 0.7rem;
  bottom: 0.7rem;
  padding: 0.28rem 0.55rem;
  border-radius: 999px;
  font-size: 0.75rem;
  color: #fff;
  background: color-mix(in srgb, #0b0c0f 72%, transparent);
  opacity: 0;
  transition: opacity 0.15s ease;
  pointer-events: none;
}

.site-shot__trigger:hover .site-shot__hint,
.site-shot__trigger:focus-visible .site-shot__hint {
  opacity: 1;
}

.site-lightbox {
  position: fixed;
  inset: 0;
  z-index: 80;
  display: grid;
  grid-template-rows: auto 1fr;
  padding: 0.85rem 0.85rem 1rem;
}

.site-lightbox__backdrop {
  position: absolute;
  inset: 0;
  border: 0;
  background: color-mix(in srgb, #05060a 82%, transparent);
  cursor: zoom-out;
}

.site-lightbox__chrome,
.site-lightbox__stage {
  position: relative;
  z-index: 1;
}

.site-lightbox__chrome {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
}

.site-lightbox__title {
  margin: 0;
  font-family: var(--site-font-display);
  font-weight: 600;
  letter-spacing: -0.02em;
}

.site-lightbox__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.site-lightbox__stage {
  min-height: 0;
  display: grid;
  place-items: center;
  overflow: auto;
}

.site-lightbox__stage--actual {
  place-items: start center;
}

.site-lightbox__img {
  display: block;
  width: auto;
  height: auto;
  max-width: min(96vw, var(--shot-w));
  max-height: calc(100vh - 6.5rem);
  border-radius: 0.75rem;
  box-shadow: 0 24px 64px color-mix(in srgb, #000 55%, transparent);
}

.site-lightbox__img--actual {
  max-width: none;
  max-height: none;
  width: var(--shot-w);
  height: auto;
}

.site-lightbox__nav {
  position: fixed;
  top: 50%;
  z-index: 2;
  width: 2.5rem;
  height: 2.5rem;
  border: 1px solid color-mix(in srgb, #fff 16%, transparent);
  border-radius: 999px;
  background: color-mix(in srgb, #0b0c0f 70%, transparent);
  color: #fff;
  font-size: 1.5rem;
  line-height: 1;
  cursor: pointer;
  transform: translateY(-50%);
}

.site-lightbox__nav--prev {
  left: 0.75rem;
}

.site-lightbox__nav--next {
  right: 0.75rem;
}

.site-lightbox__nav:hover {
  background: color-mix(in srgb, #0b0c0f 88%, transparent);
}

@media (prefers-reduced-motion: reduce) {
  .site-shot__hint {
    transition: none;
  }
}
</style>
