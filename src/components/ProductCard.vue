<script setup lang="ts">
import { RsBadge, RsButton } from '@/ui'
import { useRouter } from 'vue-router'
import {
  statusBadgeVariant,
  statusLabel,
  type ProductCardItem,
} from '../data/products'

const props = defineProps<{
  product: ProductCardItem
}>()

const router = useRouter()

function open() {
  if (!props.product.to) return
  void router.push(props.product.to)
}
</script>

<template>
  <article
    class="product-tile"
    :class="{ 'product-tile--soon': !product.to }"
  >
    <div v-if="product.shot" class="product-tile__shot">
      <img
        :src="product.shot"
        alt=""
        width="1024"
        height="640"
      />
    </div>
    <div class="product-tile__body">
      <div class="product-tile__top">
        <div class="product-tile__identity">
          <img
            v-if="product.id === 'niuma-desktop'"
            class="product-tile__icon"
            src="/brand/app-icon.svg"
            alt=""
            width="36"
            height="36"
          />
          <div
            v-else
            class="product-tile__icon product-tile__icon--placeholder"
            aria-hidden="true"
          />
          <div>
            <h3 class="product-tile__name">{{ product.name }}</h3>
            <RsBadge :variant="statusBadgeVariant[product.status]">
              {{ statusLabel[product.status] }}
            </RsBadge>
          </div>
        </div>
      </div>
      <p class="product-tile__desc">{{ product.tagline }}</p>
      <div class="product-tile__footer">
        <RsButton
          v-if="product.to"
          variant="secondary"
          size="sm"
          radius="md"
          @click="open"
        >
          了解更多
        </RsButton>
        <span v-else class="product-tile__soon">敬请期待</span>
      </div>
    </div>
  </article>
</template>

<style scoped>
.product-tile {
  display: flex;
  flex-direction: column;
  min-height: 14rem;
  padding: 0;
  overflow: hidden;
  border-radius: 1rem;
  border: 1px solid color-mix(in srgb, var(--rs-border) 90%, transparent);
  background:
    linear-gradient(165deg, color-mix(in srgb, #161922 92%, transparent), #0e1014);
  transition: border-color 0.2s ease, transform 0.2s ease, box-shadow 0.2s ease;
}

.product-tile__shot {
  aspect-ratio: 16 / 10;
  overflow: hidden;
  background: #0b0c0f;
}

.product-tile__shot img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: top left;
}

.product-tile__body {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  flex: 1;
  padding: 1.35rem 1.4rem 1.25rem;
}

.product-tile:hover {
  border-color: color-mix(in srgb, var(--site-aurora-a) 40%, var(--rs-border));
  box-shadow: 0 18px 48px color-mix(in srgb, #000 40%, transparent);
  transform: translateY(-3px);
}

.product-tile--soon {
  opacity: 0.78;
}

.product-tile__identity {
  display: flex;
  align-items: flex-start;
  gap: 0.85rem;
}

.product-tile__icon {
  border-radius: 22%;
  flex-shrink: 0;
}

.product-tile__icon--placeholder {
  width: 36px;
  height: 36px;
  border-radius: 22%;
  border: 1px dashed color-mix(in srgb, #fff 18%, transparent);
  background: color-mix(in srgb, #fff 4%, transparent);
}

.product-tile__name {
  margin: 0 0 0.45rem;
  font-family: var(--site-font-display);
  font-size: 1.2rem;
  font-weight: 650;
  letter-spacing: -0.03em;
}

.product-tile__desc {
  margin: 0;
  flex: 1;
  color: var(--rs-muted);
  line-height: 1.65;
  font-size: 0.98rem;
}

.product-tile__footer {
  display: flex;
  align-items: center;
  min-height: 2rem;
}

.product-tile__soon {
  font-size: 0.9rem;
  color: var(--rs-muted);
}

@media (prefers-reduced-motion: reduce) {
  .product-tile,
  .product-tile:hover {
    transition: none;
    transform: none;
  }
}
</style>
