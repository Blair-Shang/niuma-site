package httpapi

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/Blair-Shang/niuma-site/server/internal/config"
	"github.com/Blair-Shang/niuma-site/server/internal/ratelimit"
	"github.com/Blair-Shang/niuma-site/server/internal/store"
)

// Server 官网 HTTP 服务：下载计次 API、健康检查与内嵌静态资源。
type Server struct {
	cfg     config.Config
	store   *store.Store
	hits    *ratelimit.CoolDown
	log     *zap.Logger
	static  http.Handler
	version string
}

// New 构造官网 HTTP 服务。未调用 SetVersion 时 /healthz 的 version 为 "dev"。
func New(cfg config.Config, st *store.Store, logger *zap.Logger) *Server {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Server{
		cfg:     cfg,
		store:   st,
		hits:    ratelimit.NewCoolDown(cfg.DownloadHitCooldown),
		log:     logger,
		version: "dev",
	}
}

// SetVersion 设置发版号，写入 /healthz；空字符串忽略。
func (s *Server) SetVersion(v string) {
	v = strings.TrimSpace(v)
	if v == "" {
		return
	}
	s.version = v
}

func (s *Server) SetStatic(h http.Handler) {
	s.static = h
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("GET /api/v1/downloads/stats", s.handleStats)
	mux.HandleFunc("GET /api/v1/downloads/{platform}/hit", s.handleHit)
	mux.HandleFunc("POST /api/v1/downloads/{platform}/hit", s.handleHit)

	var root http.Handler = mux
	if s.static != nil {
		root = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" || strings.HasPrefix(r.URL.Path, "/api/") {
				mux.ServeHTTP(w, r)
				return
			}
			s.static.ServeHTTP(w, r)
		})
	}
	return s.withSecurityHeaders(s.withCORS(root))
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "version": s.version})
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.store.Stats(r.Context())
	if err != nil {
		s.log.Error("stats failed", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "stats_unavailable"})
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=15")
	writeJSON(w, http.StatusOK, stats)
}

// handleHit 302 到安装包。跳 cloud /updates/download 时由 cloud 记全局次数（含桌面自动更新）。
// 仅直链回落（windows_url）才写本站 JSON，避免和 cloud 各记一本。
// Windows 走 stable；Linux / macOS 走 preview_channel（默认 beta）。
func (s *Server) handleHit(w http.ResponseWriter, r *http.Request) {
	platform := strings.ToLower(r.PathValue("platform"))
	target, version, ok := s.cfg.DownloadURL(platform)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "unknown_platform"})
		return
	}

	ip := s.clientIP(r)
	key := platform + "|" + ip
	counted := s.hits.Allow(key)
	// 跳到 cloud /updates/download 时由 cloud 记全局次数，避免和桌面自动更新各记一本。
	if counted && !isCloudLatestDownload(target) {
		p := normalizePlatform(platform)
		if err := s.store.RecordDownload(r.Context(), p, version); err != nil {
			s.log.Error("record download failed", zap.Error(err), zap.String("platform", p))
		} else {
			s.log.Info("download hit", zap.String("platform", p), zap.String("version", version))
		}
	}

	http.Redirect(w, r, target, http.StatusFound)
}

func isCloudLatestDownload(target string) bool {
	return strings.Contains(target, "/api/v1/updates/download")
}

func normalizePlatform(p string) string {
	switch p {
	case "win", "win64", "windows":
		return "windows"
	case "mac", "osx", "darwin", "macos":
		return "macos"
	default:
		return p
	}
}

func (s *Server) clientIP(r *http.Request) string {
	remote := remoteHost(r)
	if !s.trustedProxy(remote) {
		return remote
	}
	if xri := strings.TrimSpace(r.Header.Get("X-Real-IP")); xri != "" {
		if ip := net.ParseIP(xri); ip != nil {
			return ip.String()
		}
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		candidate := strings.TrimSpace(strings.Split(xff, ",")[0])
		if ip := net.ParseIP(candidate); ip != nil {
			return ip.String()
		}
	}
	return remote
}

func remoteHost(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (s *Server) trustedProxy(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	want := parsed.String()
	for _, p := range s.cfg.TrustedProxies {
		if p == want {
			return true
		}
	}
	return false
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	allowed := map[string]struct{}{}
	for _, o := range s.cfg.CORSOrigins {
		allowed[o] = struct{}{}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			if _, ok := allowed[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			}
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
