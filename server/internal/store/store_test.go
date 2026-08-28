package store_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Blair-Shang/niuma-site/server/internal/store"
)

func TestFileRecordAndStats(t *testing.T) {
	path := filepath.Join(t.TempDir(), "download-stats.json")
	st := store.New(path)
	ctx := context.Background()

	if err := st.RecordDownload(ctx, "windows", "0.1.0"); err != nil {
		t.Fatal(err)
	}
	if err := st.RecordDownload(ctx, "windows", "0.1.0"); err != nil {
		t.Fatal(err)
	}
	if err := st.RecordDownload(ctx, "windows", "0.2.0"); err != nil {
		t.Fatal(err)
	}

	stats, err := st.Stats(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if stats.Total != 3 {
		t.Fatalf("total=%d", stats.Total)
	}
	if stats.ByPlatform["windows"] != 3 {
		t.Fatalf("byPlatform=%v", stats.ByPlatform)
	}
	if stats.ByVersion["0.1.0"] != 2 || stats.ByVersion["0.2.0"] != 1 {
		t.Fatalf("byVersion=%v", stats.ByVersion)
	}
}
