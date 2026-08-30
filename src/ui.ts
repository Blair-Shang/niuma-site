/**
 * 官网用的 @niuma/ui 薄封装：只 re-export 营销站需要的 Rs* 组件。
 * 禁止从此处导出 Monaco / Terminal / Table 等重型业务控件，避免打进静态包。
 *
 * 用法：`import { RsButton, RsCard } from '@/ui'`
 */
export {
  RsConfigProvider,
  RsButton,
  RsInput,
  RsCard,
  RsBadge,
  RsLink,
  RsIcon,
} from '@niuma/ui'
export type { RsCardVariant } from '@niuma/ui'
