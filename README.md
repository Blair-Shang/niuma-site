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

本机开发需要同级 `../niuma-ui`（`package.json` 为 `link:`）。前端环境变量见 `.env.example`。

下载链路：页面用 cloud `updates/latest` 展示版本说明；点击走本站 `/api/v1/downloads/{platform}/hit`（计次）再 302 到 cloud `/api/v1/updates/download`。Windows 为 `stable`，Linux / macOS 为 `beta`。生产同域把 `cloud.api_base` 保持为 `/niuma/cloud` 即可。

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

生产反代（与现有 nginx 互不干扰）：

1. 备份 `/data/nginx/conf/nginx.conf`
2. 把原 `server { minihub... }` 挪到 `conf.d/minihub.net.cn.conf`（样例 [config/conf.d/minihub.net.cn.conf.example](./config/conf.d/minihub.net.cn.conf.example)），`routes/*.conf` 不用动
3. 主配置改成只 `include conf.d/*.conf;`（样例 [config/nginx.conf.example](./config/nginx.conf.example)）
4. 放入 [config/conf.d/niuma-site.conf](./config/conf.d/niuma-site.conf) 和 niuma007 证书
5. `nginx -t` 后 reload

## GitHub 打包

推送 `v*` 标签（须与 `package.json` 的 `version` 一致）或手动运行 Actions **Pack and Release**。

产物：

- `niuma-site-<version>-linux-amd64.tar.gz`（解压后目录为 `niuma-site/`）
- `niuma-site-<version>-windows-amd64.zip`
- `SHA256SUMS.txt`（GNU `sha256sum` 文本格式）

Release 正文从 `CHANGELOG.md` 对应 `## [x.y.z]` 段生成；tag、`package.json` 的 `version` 与 changelog 版本号必须一致。发版前把 `[Unreleased]` 条目迁入新版本段。

CI / 打包从 npm 安装 `niuma-ui@latest`（可用变量 `NIUMA_UI_VERSION` 钉死如 `1.2.0`，或手动 Pack 时填写 `niuma_ui_version`）。不 checkout 兄弟 Git 仓。没有 GitHub 标签 `latest`。
