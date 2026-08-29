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
	DownloadHitCooldown    time.Duration

	// CloudAPIBase 官网 hit 优先 302 的 cloud 根（含 /niuma/cloud，不含 /api/v1）。
	// 生产同域用相对路径 /niuma/cloud；设为 off 或 - 时只用 DownloadWindowsURL。
	CloudAPIBase        string
	CloudProduct        string
	CloudChannel        string
	CloudPreviewChannel string
	CloudWindowsArch    string
	CloudLinuxArch      string
	CloudMacOSArch      string

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
		HitCooldownSec int    `yaml:"hit_cooldown_sec"`
	} `yaml:"download"`

	Cloud struct {
		APIBase        string `yaml:"api_base"`
		Product        string `yaml:"product"`
		Channel        string `yaml:"channel"`
		PreviewChannel string `yaml:"preview_channel"`
		WindowsArch    string `yaml:"windows_arch"`
		LinuxArch      string `yaml:"linux_arch"`
		MacOSArch      string `yaml:"macos_arch"`
	} `yaml:"cloud"`

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
		CloudAPIBase:           normalizeCloudAPIBase(fc.Cloud.APIBase),
		CloudProduct:           firstNonEmpty(strings.TrimSpace(fc.Cloud.Product), "niuma"),
		CloudChannel:           firstNonEmpty(strings.TrimSpace(fc.Cloud.Channel), "stable"),
		CloudPreviewChannel:    firstNonEmpty(strings.TrimSpace(fc.Cloud.PreviewChannel), "beta"),
		CloudWindowsArch:       firstNonEmpty(strings.TrimSpace(fc.Cloud.WindowsArch), "x64"),
		CloudLinuxArch:         firstNonEmpty(strings.TrimSpace(fc.Cloud.LinuxArch), "x64"),
		CloudMacOSArch:         firstNonEmpty(strings.TrimSpace(fc.Cloud.MacOSArch), "arm64"),
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
	if cfg.CloudAPIBase != "" {
		for _, plat := range []string{"windows", "linux", "macos"} {
			if loc := cfg.CloudLatestDownloadURL(plat); loc == "" {
				return Config{}, fmt.Errorf("cloud.api_base: invalid latest download location for %s", plat)
			}
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
cloud:
  api_base: "/niuma/cloud"
  product: niuma
  channel: stable
  preview_channel: beta
  windows_arch: x64
  linux_arch: x64
  macos_arch: arm64
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

// DownloadURL 返回 hit 302 目标。
// Windows 走 stable（可回落 windows_url）；Linux / macOS 走 preview_channel（默认 beta）。
func (c Config) DownloadURL(platform string) (target, version string, ok bool) {
	switch strings.ToLower(platform) {
	case "windows", "win", "win64":
		if loc := c.CloudLatestDownloadURL("windows"); loc != "" {
			return loc, "latest", true
		}
		if err := ValidateHTTPSDownloadURL(c.DownloadWindowsURL); err != nil {
			return "", "", false
		}
		return c.DownloadWindowsURL, c.DownloadWindowsVersion, true
	case "linux":
		if loc := c.CloudLatestDownloadURL("linux"); loc != "" {
			return loc, "latest", true
		}
	case "macos", "mac", "osx", "darwin":
		if loc := c.CloudLatestDownloadURL("macos"); loc != "" {
			return loc, "latest", true
		}
	}
	return "", "", false
}

// CloudLatestDownloadURL 拼出 cloud「点击当下解析最新 published」的稳定地址。
func (c Config) CloudLatestDownloadURL(platform string) string {
	base := strings.TrimRight(strings.TrimSpace(c.CloudAPIBase), "/")
	if base == "" {
		return ""
	}
	channel, arch, ok := c.platformDownloadDim(platform)
	if !ok {
		return ""
	}
	product := firstNonEmpty(c.CloudProduct, "niuma")
	q := url.Values{}
	q.Set("product", product)
	q.Set("channel", channel)
	q.Set("platform", platform)
	q.Set("arch", arch)
	loc := base + "/api/v1/updates/download?" + q.Encode()
	if err := ValidateDownloadLocation(loc); err != nil {
		return ""
	}
	return loc
}

func (c Config) platformDownloadDim(platform string) (channel, arch string, ok bool) {
	switch strings.ToLower(strings.TrimSpace(platform)) {
	case "windows":
		return firstNonEmpty(c.CloudChannel, "stable"), firstNonEmpty(c.CloudWindowsArch, "x64"), true
	case "linux":
		return firstNonEmpty(c.CloudPreviewChannel, "beta"), firstNonEmpty(c.CloudLinuxArch, "x64"), true
	case "macos":
		return firstNonEmpty(c.CloudPreviewChannel, "beta"), firstNonEmpty(c.CloudMacOSArch, "arm64"), true
	default:
		return "", "", false
	}
}

func normalizeCloudAPIBase(raw string) string {
	raw = strings.TrimSpace(raw)
	switch strings.ToLower(raw) {
	case "off", "-":
		return ""
	case "":
		return "/niuma/cloud"
	default:
		return strings.TrimRight(raw, "/")
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

// ValidateDownloadLocation 允许：cloud 最新包相对路径、https，或本机回环 http（开发）。
func ValidateDownloadLocation(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Errorf("empty")
	}
	if strings.HasPrefix(raw, "/") {
		if strings.HasPrefix(raw, "//") || strings.Contains(raw, "\\") || strings.Contains(raw, "@") {
			return fmt.Errorf("invalid relative download path")
		}
		pathOnly := raw
		if i := strings.IndexByte(raw, '?'); i >= 0 {
			pathOnly = raw[:i]
		}
		if pathOnly != "/niuma/cloud/api/v1/updates/download" {
			return fmt.Errorf("relative path must be cloud latest download")
		}
		return nil
	}
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if u.User != nil || u.Host == "" {
		return fmt.Errorf("invalid url")
	}
	if !strings.HasSuffix(u.Path, "/api/v1/updates/download") {
		return fmt.Errorf("path must end with /api/v1/updates/download")
	}
	if u.Scheme == "https" {
		return nil
	}
	if u.Scheme == "http" && isLoopbackHost(u.Hostname()) {
		return nil
	}
	return fmt.Errorf("must be https, loopback http, or /niuma/cloud/api/v1/updates/download")
}

func isLoopbackHost(host string) bool {
	switch strings.ToLower(strings.TrimSpace(host)) {
	case "127.0.0.1", "localhost", "localhost.", "::1":
		return true
	default:
		return false
	}
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
