package store

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// Store 用本地 JSON 元数据文件累计下载次数（无数据库）。
type Store struct {
	path string
	mu   sync.Mutex
}

type DownloadStats struct {
	Total      int64            `json:"total"`
	ByPlatform map[string]int64 `json:"byPlatform"`
	ByVersion  map[string]int64 `json:"byVersion,omitempty"`
}

func New(path string) *Store {
	return &Store{path: path}
}

func (s *Store) RecordDownload(_ context.Context, platform, version string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stats, err := s.readUnlocked()
	if err != nil {
		return err
	}
	if stats.ByPlatform == nil {
		stats.ByPlatform = map[string]int64{}
	}
	if stats.ByVersion == nil {
		stats.ByVersion = map[string]int64{}
	}
	stats.Total++
	stats.ByPlatform[platform]++
	if version != "" {
		stats.ByVersion[version]++
	}
	return s.writeUnlocked(stats)
}

func (s *Store) Stats(_ context.Context) (DownloadStats, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.readUnlocked()
}

func (s *Store) readUnlocked() (DownloadStats, error) {
	out := DownloadStats{
		ByPlatform: map[string]int64{},
		ByVersion:  map[string]int64{},
	}
	raw, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return out, nil
		}
		return out, err
	}
	if len(raw) == 0 {
		return out, nil
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return DownloadStats{}, err
	}
	if out.ByPlatform == nil {
		out.ByPlatform = map[string]int64{}
	}
	if out.ByVersion == nil {
		out.ByVersion = map[string]int64{}
	}
	return out, nil
}

func (s *Store) writeUnlocked(stats DownloadStats) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
