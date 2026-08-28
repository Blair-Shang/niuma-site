# niuma-site

[niuma007.com](https://niuma007.com) 官方网站：Vue 前端 + Go API（静态资源内嵌进 Go 二进制）。

## 开发

```bash
pnpm install
Copy-Item .env.example .env
# 后端配置：首次启动会从 config/app.yaml.example 生成 config/app.yaml

pnpm dev:server   # :8080 API
pnpm dev          # :5173 前端，/api 代理到 8080
```

需要同级 `../niuma-ui`。前端环境变量见 `.env.example`。

## 生产构建与启停

Windows：

```powershell
powershell -ExecutionPolicy Bypass -File script/build.ps1
powershell -ExecutionPolicy Bypass -File script/start.ps1
# stop / restart 同理
```

Linux：

```bash
bash script/build.sh
bash script/start.sh
```

单进程访问：`http://127.0.0.1:8080/`（页面 + `/api`）。脚本说明见 [script/README.md](./script/README.md)。

## GitHub 打包

推送 `v*` 标签（须与 `package.json` 的 `version` 一致）或手动运行 Actions **Pack and Release**。

产物：

- `niuma-site-<version>-linux-amd64.tar.gz`（内含同名目录、`CHANGELOG.md`、`VERSION`）
- `niuma-site-<version>-windows-amd64.zip`
- `SHA256SUMS.txt`（GNU `sha256sum` 文本格式）

Release 正文从 `CHANGELOG.md` 对应 `## [x.y.z]` 段生成；tag、`package.json` 的 `version` 与 changelog 版本号必须一致。发版前把 `[Unreleased]` 条目迁入新版本段。

CI 会 checkout 同组织的 `niuma-ui` 作为兄弟目录（`link:../niuma-ui`）。若 `niuma-ui` 为私有仓库，在本仓配置 secret `NIUMA_UI_TOKEN`；可用 repository variable `NIUMA_UI_REF` 覆盖默认组件库版本。
