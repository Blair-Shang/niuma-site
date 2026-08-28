<script setup lang="ts">
import { RsConfigProvider } from '@/ui'
import SiteLayout from './layouts/SiteLayout.vue'
</script>

<template>
  <RsConfigProvider theme="dark" locale="zh-CN" control-size="md">
    <SiteLayout />
  </RsConfigProvider>
</template>
