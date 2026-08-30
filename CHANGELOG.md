# 变更日志

本文件记录 niuma-site 的重要变更。

格式参考 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，版本遵循 [Semantic Versioning](https://semver.org/lang/zh-CN/)。

发版时把 `[Unreleased]` 下的条目移入新的 `## [x.y.z] - YYYY-MM-DD` 段；GitHub Release 正文由该段生成。`package.json` 的 `version`、git tag（`v*`）与本文件版本号必须一致。公开基线为 **1.0.0**。

## [Unreleased]

## [1.1.4] - 2026-08-30

### 修复

- 官网 CSS 重新 `@import 'tailwindcss'`，让扫描根在本仓库；打包关闭 `cssCodeSplit`，组件样式与 token 打进同一份 CSS。依赖 `niuma-ui` 1.2.3 起把旁路组件 CSS 与 `@source` 一并修好。

## [1.1.3] - 2026-08-30

### 变更

- `src/ui.ts` 从 `@niuma/ui` 具名导入；`niumaUiHost` 让本机 `pnpm dev` 联调源码，打包走 npm `dist`。Tailwind 只通过 `niuma-ui/styles.css` 透传，官网 CSS 不再自己 `@import`。

## [1.1.2] - 2026-08-30

### 修复

- 官网自行 `@import 'tailwindcss'`（Preflight），不再依赖本机 `niuma-ui` 源码样式。流水线用的 npm `styles.css` 不含 Tailwind，线上盒模型 / 标题 / 按钮不再和 `pnpm dev` 对不齐。

## [1.1.1] - 2026-08-30

### 变更

- 从 `@niuma/ui` 包入口具名导入组件，不再依赖 `@niuma-ui-src` 源码别名。
- CI / 打包从 npm 安装 `niuma-ui@latest`（可用 `NIUMA_UI_VERSION` 钉死），不再 checkout 兄弟 Git 仓。本机仍用 `link:../niuma-ui` 联调。

## [1.1.0] - 2026-08-30

### 新增

- 首页、产品页与产品卡片展示桌面端真实工作台截图（运维监控、数据库），替换 CSS 占位窗。
- 产品截图支持点击放大、1:1 原图像素，以及新标签打开原图。

### 变更

- Nginx 将 `/niuma/` 反代到 niuma-cloud；官网默认 Cloud API 基址为 `/niuma/cloud`。
- 下载 hit 优先 302 到 niuma-cloud `GET /api/v1/updates/download`：点击当下解析最新 published，不再依赖手改 `download.windows_url`。
- `config/app.yaml` 增加 `cloud.api_base`（默认 `/niuma/cloud`）；`windows_url` 仅在 `api_base: off` 时作为应急直链。
- 下载页开放 Linux（x64）与 macOS（arm64）：走 `preview_channel`（默认 **beta**）。Windows 仍为 `stable`。无对应 published 时按钮停用。
- 下载页展示各平台 SHA-256（来自 `updates/latest`），点击复制完整哈希。
- 下载页导语改为对外表述，去掉渠道 / published 等内部用语。
- Nginx：对 `/niuma/cloud/api/v1/updates/download` 按 IP 限请求，对 `/updates/files/` 限并发与带宽。
- CI / 打包默认 checkout 同组织 `niuma-ui` 的 `main`，可用 `NIUMA_UI_REF` 钉死。npm 消费用 `niuma-ui@latest`。

### 说明

- 下载页版本与更新说明仍来自 `updates/latest`（展示）；真正下包走 hit → cloud latest download，避免标签页开着时下到旧包。
- Linux / macOS 尚未稳定，官网与流水线按 beta 登记；稳定后再切到 `channel: stable`。

## [1.0.3] - 2026-08-29

### 变更

- GitHub Actions 升级到原生 Node 24 运行时（`checkout`/`setup-node`/`setup-go`/`pnpm/action-setup`/`upload-artifact`/`download-artifact`），消除 Node 20 弃用警告。
- 主导航「产品」改为路径 `/products`，不再使用 `/#products` 锚点。

## [1.0.2] - 2026-08-29

### 变更

- 发版归档文件名仍带版本号；解压根目录为 `niuma-site/`。

## [1.0.1] - 2026-08-29

### 新增

- `config/conf.d/niuma-site.conf`：独立 `server`，反代本站 `127.0.0.1:8080`；证书路径对齐 `/data/nginx/ssl/`。
- `config/nginx.conf.example`：主配置只 `include conf.d/*.conf;`，站点互不嵌套。
- `config/conf.d/minihub.net.cn.conf.example`：从原 `nginx.conf` 迁出既有站点的样例。

### 修复

- 去掉 `index.html` 重复 `</head>`；根目录忽略改为 `/data/`，避免 `src/data` 进不了仓库导致 CI 缺产品列表。

## [1.0.0] - 2026-08-29

### 新增

- 官网：首页、产品、下载、反馈、关于；页脚悬挂 ICP 备案号。
- Vue 前端嵌入 Go 单二进制（`run/niuma-site`），`/healthz` 与 `/api/v1/downloads/*` 同进程。
- 下载计次后 302 到 `config/app.yaml` 的 `windows_url`；版本说明与反馈走 niuma-cloud。
- GitHub Actions：CI + Pack and Release；产物含 changelog、VERSION、SHA256SUMS。

### 安全

- 响应增加 CSP / `X-Frame-Options` / `nosniff` 等安全头；下载 302 仅允许 `https` URL。
- 默认只采信来自 `trusted_proxies` 的 `X-Forwarded-For`，避免伪造 IP 绕过下载冷却。
- 生产默认监听 `127.0.0.1:8080`；去掉 Google Fonts 外链（避免第三方探测与国内不可达）。

[Unreleased]: https://github.com/Blair-Shang/niuma-site/compare/v1.1.1...HEAD
[1.1.1]: https://github.com/Blair-Shang/niuma-site/releases/tag/v1.1.1
[1.1.0]: https://github.com/Blair-Shang/niuma-site/releases/tag/v1.1.0
[1.0.3]: https://github.com/Blair-Shang/niuma-site/releases/tag/v1.0.3
[1.0.2]: https://github.com/Blair-Shang/niuma-site/releases/tag/v1.0.2
[1.0.1]: https://github.com/Blair-Shang/niuma-site/releases/tag/v1.0.1
[1.0.0]: https://github.com/Blair-Shang/niuma-site/releases/tag/v1.0.0
