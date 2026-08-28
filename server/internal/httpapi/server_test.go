package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Blair-Shang/niuma-site/server/internal/config"
	"github.com/Blair-Shang/niuma-site/server/internal/httpapi"
	"github.com/Blair-Shang/niuma-site/server/internal/store"
)

func TestHealthVersion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "download-stats.json")
	srv := httpapi.New(config.Config{}, store.New(path), nil)
	srv.SetVersion("1.0.0")
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" || body["version"] != "1.0.0" {
		t.Fatalf("body=%v", body)
	}
}

func TestHitRedirectAndStats(t *testing.T) {
	path := filepath.Join(t.TempDir(), "download-stats.json")
	cfg := config.Config{
		DownloadWindowsURL:     "https://example.com/setup.exe",
		DownloadWindowsVersion: "0.1.0",
		DownloadHitCooldown:    time.Second,
		CORSOrigins:            []string{"http://localhost:5173"},
	}
	srv := httpapi.New(cfg, store.New(path), nil)
	h := srv.Handler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/downloads/windows/hit", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("status=%d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "https://example.com/setup.exe" {
		t.Fatalf("location=%s", loc)
	}

	sreq := httptest.NewRequest(http.MethodGet, "/api/v1/downloads/stats", nil)
	srec := httptest.NewRecorder()
	h.ServeHTTP(srec, sreq)
	if srec.Code != http.StatusOK {
		t.Fatalf("stats status=%d body=%s", srec.Code, srec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(srec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["total"].(float64) != 1 {
		t.Fatalf("body=%v", body)
	}
}

func TestSecurityHeaders(t *testing.T) {
	path := filepath.Join(t.TempDir(), "download-stats.json")
	srv := httpapi.New(config.Config{}, store.New(path), nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	h := rec.Header()
	if h.Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("nosniff=%q", h.Get("X-Content-Type-Options"))
	}
	if h.Get("X-Frame-Options") != "DENY" {
		t.Fatalf("frame=%q", h.Get("X-Frame-Options"))
	}
	if !strings.Contains(h.Get("Content-Security-Policy"), "frame-ancestors 'none'") {
		t.Fatalf("csp=%q", h.Get("Content-Security-Policy"))
	}
}

func TestHitRejectsNonHTTPSDownloadURL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "download-stats.json")
	srv := httpapi.New(config.Config{
		DownloadWindowsURL:  "http://example.com/setup.exe",
		DownloadHitCooldown: time.Second,
	}, store.New(path), nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/downloads/windows/hit", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestHitIgnoresSpoofedForwardedFor(t *testing.T) {
	path := filepath.Join(t.TempDir(), "download-stats.json")
	srv := httpapi.New(config.Config{
		DownloadWindowsURL:  "https://example.com/setup.exe",
		DownloadHitCooldown: time.Minute,
	}, store.New(path), nil)
	h := srv.Handler()

	hit := func(xff string) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/downloads/windows/hit", nil)
		req.RemoteAddr = "203.0.113.10:54321"
		req.Header.Set("X-Forwarded-For", xff)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusFound {
			t.Fatalf("status=%d xff=%s", rec.Code, xff)
		}
	}
	hit("198.51.100.1")
	hit("198.51.100.2")

	srec := httptest.NewRecorder()
	h.ServeHTTP(srec, httptest.NewRequest(http.MethodGet, "/api/v1/downloads/stats", nil))
	var body map[string]any
	if err := json.Unmarshal(srec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["total"].(float64) != 1 {
		t.Fatalf("spoofed XFF counted separately: %v", body)
	}
}
