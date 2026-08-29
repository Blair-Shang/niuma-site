package config

import "testing"

func TestCloudLatestDownloadURL(t *testing.T) {
	c := Config{
		CloudAPIBase:        "/niuma/cloud",
		CloudProduct:        "niuma",
		CloudChannel:        "stable",
		CloudPreviewChannel: "beta",
		CloudWindowsArch:    "x64",
		CloudLinuxArch:      "x64",
		CloudMacOSArch:      "arm64",
	}
	got := c.CloudLatestDownloadURL("windows")
	want := "/niuma/cloud/api/v1/updates/download?arch=x64&channel=stable&platform=windows&product=niuma"
	if got != want {
		t.Fatalf("windows=%s", got)
	}
	linux := c.CloudLatestDownloadURL("linux")
	wantLinux := "/niuma/cloud/api/v1/updates/download?arch=x64&channel=beta&platform=linux&product=niuma"
	if linux != wantLinux {
		t.Fatalf("linux=%s", linux)
	}
	mac := c.CloudLatestDownloadURL("macos")
	wantMac := "/niuma/cloud/api/v1/updates/download?arch=arm64&channel=beta&platform=macos&product=niuma"
	if mac != wantMac {
		t.Fatalf("macos=%s", mac)
	}

	target, version, ok := c.DownloadURL("windows")
	if !ok || version != "latest" || target != want {
		t.Fatalf("download url=%s version=%s ok=%v", target, version, ok)
	}
	if loc, _, ok := c.DownloadURL("linux"); !ok || loc != wantLinux {
		t.Fatalf("linux hit=%s ok=%v", loc, ok)
	}
	if loc, _, ok := c.DownloadURL("darwin"); !ok || loc != wantMac {
		t.Fatalf("darwin hit=%s ok=%v", loc, ok)
	}
}

func TestDownloadURLFallsBackToWindowsURL(t *testing.T) {
	c := Config{DownloadWindowsURL: "https://cdn.example.com/setup.exe", DownloadWindowsVersion: "1.0.0"}
	target, version, ok := c.DownloadURL("windows")
	if !ok || target != c.DownloadWindowsURL || version != "1.0.0" {
		t.Fatalf("fallback target=%s version=%s ok=%v", target, version, ok)
	}
}

func TestValidateDownloadLocation(t *testing.T) {
	ok := []string{
		"/niuma/cloud/api/v1/updates/download?platform=windows&arch=x64",
		"https://www.niuma007.com/niuma/cloud/api/v1/updates/download",
		"http://127.0.0.1:8090/niuma/cloud/api/v1/updates/download",
	}
	for _, raw := range ok {
		if err := ValidateDownloadLocation(raw); err != nil {
			t.Fatalf("want ok %s: %v", raw, err)
		}
	}
	bad := []string{
		"",
		"//evil.example/niuma/cloud/api/v1/updates/download",
		"/api/v1/updates/download",
		"http://example.com/niuma/cloud/api/v1/updates/download",
		"https://www.niuma007.com/other",
	}
	for _, raw := range bad {
		if err := ValidateDownloadLocation(raw); err == nil {
			t.Fatalf("want reject %s", raw)
		}
	}
}

func TestNormalizeCloudAPIBase(t *testing.T) {
	if got := normalizeCloudAPIBase(""); got != "/niuma/cloud" {
		t.Fatalf("default=%s", got)
	}
	if got := normalizeCloudAPIBase("off"); got != "" {
		t.Fatalf("off=%s", got)
	}
	if got := normalizeCloudAPIBase("https://www.niuma007.com/niuma/cloud/"); got != "https://www.niuma007.com/niuma/cloud" {
		t.Fatalf("abs=%s", got)
	}
}
