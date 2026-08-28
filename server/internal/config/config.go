package config

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPAddr string

	ConfigDir string
	LogDir    string
	DataDir   string

	// DownloadStatsFile 累计下载元数据 JSON（无数据库）
	DownloadStatsFile string

	CORSOrigins []string

	DownloadWindowsURL     string
	DownloadWindowsVersion string
	DownloadHitCooldown     time.Duration

	// TrustedProxies 反代 IP；仅这些来源才采信 X-Real-IP / X-Forwarded-For。
	TrustedProxies []string

	FeedbackWebhookURL string

	LogLevel      string
	LogToConsole  bool
	LogMaxSizeMB  int
	LogMaxBackups int
	LogMaxAgeDays int
}

// fileConfig 对应 config/app.yaml。
type fileConfig struct {
	HTTPAddr string `yaml:"http_addr"`

	DataDir string `yaml:"data_dir"`

	Download struct {
		StatsFile      string `yaml:"stats_file"`
		WindowsURL     string `yaml:"windows_url"`
		WindowsVersion string `yaml:"windows_version"`
		HitCooldownSec  int    `yaml:"hit_cooldown_sec"`
	} `yaml:"download"`

	CORSOrigins []string `yaml:"cors_origins"`

	TrustedProxies []string `yaml:"trusted_proxies"`

	Feedback struct {
		WebhookURL string `yaml:"webhook_url"`
	} `yaml:"feedback"`

	Log struct {
		Level      string `yaml:"level"`
		ToConsole  *bool  `yaml:"to_console"`
		Dir        string `yaml:"dir"`
		MaxSizeMB  int    `yaml:"max_size_mb"`
		MaxBackups int    `yaml:"max_backups"`
		MaxAgeDays int    `yaml:"max_age_days"`
	} `yaml:"log"`
}

// Load 确保目录存在，再读取 config/app.yaml。
// 仅 CONFIG_DIR 可用环境变量覆盖配置目录。
func Load() (Config, error) {
	cfgDir := getenv("CONFIG_DIR", "./config")
	if err := ensureRuntimeDirs(cfgDir); err != nil {
		return Config{}, err
	}
	if err := ensureConfigFile(cfgDir); err != nil {
		return Config{}, err
	}

	raw, err := os.ReadFile(filepath.Join(cfgDir, "app.yaml"))
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var fc fileConfig
	if err := yaml.Unmarshal(raw, &fc); err != nil {
		return Config{}, fmt.Errorf("parse app.yaml: %w", err)
	}

	logDir := firstNonEmpty(fc.Log.Dir, "./logs")
	dataDir := firstNonEmpty(fc.DataDir, "./data")
	if err := ensureRuntimeDirs(logDir, dataDir); err != nil {
		return Config{}, err
	}

	logToConsole := true
	if fc.Log.ToConsole != nil {
		logToConsole = *fc.Log.ToConsole
	}

	cooldownSec := fc.Download.HitCooldownSec
	if cooldownSec < 0 {
		return Config{}, fmt.Errorf("download.hit_cooldown_sec invalid")
	}
	if cooldownSec == 0 {
		cooldownSec = 30
	}

	cfg := Config{
		HTTPAddr:               firstNonEmpty(fc.HTTPAddr, "127.0.0.1:8080"),
		ConfigDir:              cfgDir,
		LogDir:                 logDir,
		DataDir:                dataDir,
		DownloadStatsFile:      firstNonEmpty(fc.Download.StatsFile, filepath.Join(dataDir, "download-stats.json")),
		DownloadWindowsURL:     strings.TrimSpace(fc.Download.WindowsURL),
		DownloadWindowsVersion: strings.TrimSpace(fc.Download.WindowsVersion),
		DownloadHitCooldown:    time.Duration(cooldownSec) * time.Second,
		FeedbackWebhookURL:     strings.TrimSpace(fc.Feedback.WebhookURL),
		LogLevel:               firstNonEmpty(fc.Log.Level, "info"),
		LogToConsole:           logToConsole,
		LogMaxSizeMB:           intOr(fc.Log.MaxSizeMB, 100),
		LogMaxBackups:          intOr(fc.Log.MaxBackups, 10),
		LogMaxAgeDays:          intOr(fc.Log.MaxAgeDays, 30),
	}

	if len(fc.CORSOrigins) == 0 {
		cfg.CORSOrigins = []string{
			"http://localhost:5173",
			"https://niuma007.com",
			"https://www.niuma007.com",
		}
	} else {
		for _, o := range fc.CORSOrigins {
			o = strings.TrimSpace(o)
			if o != "" {
				cfg.CORSOrigins = append(cfg.CORSOrigins, o)
			}
		}
	}

	if strings.TrimSpace(cfg.DownloadWindowsURL) != "" {
		if err := ValidateHTTPSDownloadURL(cfg.DownloadWindowsURL); err != nil {
			return Config{}, fmt.Errorf("download.windows_url: %w", err)
		}
	}

	for _, raw := range fc.TrustedProxies {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		ip := net.ParseIP(raw)
		if ip == nil {
			return Config{}, fmt.Errorf("trusted_proxies invalid IP: %s", raw)
		}
		cfg.TrustedProxies = append(cfg.TrustedProxies, ip.String())
	}

	if strings.TrimSpace(cfg.DownloadStatsFile) == "" {
		return Config{}, fmt.Errorf("download.stats_file is required")
	}
	return cfg, nil
}

func ensureRuntimeDirs(dirs ...string) error {
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", d, err)
		}
	}
	return nil
}

func ensureConfigFile(cfgDir string) error {
	dst := filepath.Join(cfgDir, "app.yaml")
	if _, err := os.Stat(dst); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	candidates := []string{
		filepath.Join(cfgDir, "app.yaml.example"),
		"config/app.yaml.example",
	}
	var src string
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			src = c
			break
		}
	}
	if src == "" {
		return os.WriteFile(dst, []byte(defaultYAML()), 0o644)
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func defaultYAML() string {
	return `http_addr: "127.0.0.1:8080"
data_dir: ./data
download:
  stats_file: ./data/download-stats.json
  windows_url: ""
  windows_version: ""
  hit_cooldown_sec: 30
trusted_proxies:
  - 127.0.0.1
  - ::1
`
}

func (c Config) DownloadURL(platform string) (url, version string, ok bool) {
	switch strings.ToLower(platform) {
	case "windows", "win", "win64":
		if err := ValidateHTTPSDownloadURL(c.DownloadWindowsURL); err != nil {
			return "", "", false
		}
		return c.DownloadWindowsURL, c.DownloadWindowsVersion, true
	default:
		return "", "", false
	}
}

// ValidateHTTPSDownloadURL 要求 https 且无 userinfo，避免误配成开放重定向。
func ValidateHTTPSDownloadURL(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Errorf("empty")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if u.Scheme != "https" || u.Host == "" || u.User != nil {
		return fmt.Errorf("must be https URL without credentials")
	}
	return nil
}

func getenv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func intOr(v, fallback int) int {
	if v == 0 {
		return fallback
	}
	return v
}
