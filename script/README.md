# 运维脚本

仓库根目录执行。Windows 用 `.ps1`，Linux 用 `.sh`。

| 脚本 | 作用 |
|------|------|
| `build` | `pnpm` 构建前端到 `server/internal/web/dist`，再 `go build` 打出单二进制（**内嵌静态资源**） |
| `pack-release` | 校验 CHANGELOG 版本段，打出 Linux/Windows amd64 发版归档到 `output/` |
| `start` | 后台启动 `run/niuma-site[.exe]`，写 pid / 日志 |
| `stop` | 按 pid 停止 |
| `restart` | stop + start |

## Windows

```powershell
powershell -ExecutionPolicy Bypass -File script/build.ps1
powershell -ExecutionPolicy Bypass -File script/start.ps1
powershell -ExecutionPolicy Bypass -File script/stop.ps1
powershell -ExecutionPolicy Bypass -File script/restart.ps1
```

## Linux

```bash
chmod +x script/*.sh
bash script/build.sh
bash script/start.sh
bash script/stop.sh
bash script/restart.sh
```

## 产物

| 路径 | 说明 |
|------|------|
| `run/niuma-site.exe` / `run/niuma-site` | 可执行文件（内嵌前端） |
| `run/niuma-site.pid` | 进程号 |
| `config/` | 启动时自动创建；`app.yaml` 首次从 `app.yaml.example` 生成 |
| `logs/` | 启动时自动创建；业务日志 `niuma-site.log`（zap + 滚动） |
| `data/` | 启动时自动创建；`download-stats.json` 累计下载 |

默认监听 `http://127.0.0.1:8080/`（`http_addr`，见 `config/app.yaml.example`）。

日志使用 **uber-go/zap** + lumberjack 滚动，文件位于 `logs/niuma-site.log`。

开发联调仍可用 `pnpm dev` + `pnpm dev:server`（Vite 热更新，不走内嵌）。
