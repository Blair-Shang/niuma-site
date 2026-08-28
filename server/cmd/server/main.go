package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Blair-Shang/niuma-site/server/internal/config"
	"github.com/Blair-Shang/niuma-site/server/internal/httpapi"
	"github.com/Blair-Shang/niuma-site/server/internal/logging"
	"github.com/Blair-Shang/niuma-site/server/internal/store"
	"github.com/Blair-Shang/niuma-site/server/internal/web"
)

// version 由打包脚本 -ldflags "-X main.version=x.y.z" 注入；本地构建为 dev。
var version = "dev"

func main() {
	if wantsVersion(os.Args[1:]) {
		fmt.Println(version)
		return
	}

	// go run -C server / 从 server 目录启动时，切回仓库根以便加载 config/logs/data
	chdirRepoRoot()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger, err := logging.New(logging.Options{
		Dir:        cfg.LogDir,
		Level:      cfg.LogLevel,
		ToConsole:  cfg.LogToConsole,
		MaxSizeMB:  cfg.LogMaxSizeMB,
		MaxBackups: cfg.LogMaxBackups,
		MaxAgeDays: cfg.LogMaxAgeDays,
	})
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	st := store.New(cfg.DownloadStatsFile)
	api := httpapi.New(cfg, st, logger)
	api.SetVersion(version)

	static, err := web.Handler()
	if err != nil {
		logger.Fatal("embed web", zap.Error(err))
	}
	api.SetStatic(static)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           api.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    16 << 10,
	}

	go func() {
		logger.Info("server listening",
			zap.String("version", version),
			zap.String("addr", cfg.HTTPAddr),
			zap.String("statsFile", cfg.DownloadStatsFile),
			zap.String("configDir", cfg.ConfigDir),
			zap.String("logDir", cfg.LogDir),
			zap.String("dataDir", cfg.DataDir),
			zap.String("web", "embedded"),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Warn("shutdown", zap.Error(err))
	} else {
		logger.Info("shutdown complete")
	}
}

func wantsVersion(args []string) bool {
	for _, a := range args {
		if a == "--version" || a == "-version" || a == "-V" {
			return true
		}
	}
	return false
}

func chdirRepoRoot() {
	if _, err := os.Stat("server/go.mod"); err == nil {
		return
	}
	if _, err := os.Stat("go.mod"); err == nil {
		_ = os.Chdir("..")
	}
}
