/**
 * 官网用的 @niuma/ui 薄封装：只 re-export 营销站需要的 Rs* 组件。
 * 禁止从此处导出 Monaco / Terminal / Table 等重型业务控件，避免打进静态包。
 *
 * 用法：`import { RsButton, RsCard } from '@/ui'`
 */
export { default as RsConfigProvider } from '@niuma-ui-src/components/RsConfigProvider.vue'
export { default as RsButton } from '@niuma-ui-src/components/RsButton.vue'
export { default as RsInput } from '@niuma-ui-src/components/RsInput.vue'
export { default as RsCard } from '@niuma-ui-src/components/RsCard.vue'
export type { RsCardVariant } from '@niuma-ui-src/components/RsCard.vue'
export { default as RsBadge } from '@niuma-ui-src/components/RsBadge.vue'
export { default as RsLink } from '@niuma-ui-src/components/RsLink.vue'
export { default as RsIcon } from '@niuma-ui-src/components/RsIcon.vue'
