# 变更日志

本文件记录 niuma-site 的重要变更。

格式参考 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，版本遵循 [Semantic Versioning](https://semver.org/lang/zh-CN/)。

发版时把 `[Unreleased]` 下的条目移入新的 `## [x.y.z] - YYYY-MM-DD` 段；GitHub Release 正文由该段生成。`package.json` 的 `version`、git tag（`v*`）与本文件版本号必须一致。公开基线为 **1.0.0**。

## [Unreleased]

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

[Unreleased]: https://github.com/Blair-Shang/niuma-site/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/Blair-Shang/niuma-site/releases/tag/v1.0.0
