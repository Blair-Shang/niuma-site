package ratelimit

import (
	"sync"
	"time"
)

// CoolDown 按 key 冷却，用于下载 hit 防刷。
type CoolDown struct {
	mu   sync.Mutex
	last map[string]time.Time
	ttl  time.Duration
}

func NewCoolDown(ttl time.Duration) *CoolDown {
	return &CoolDown{
		last: make(map[string]time.Time),
		ttl:  ttl,
	}
}

// Allow 若距上次不足 ttl 返回 false。
func (c *CoolDown) Allow(key string) bool {
	if c.ttl <= 0 {
		return true
	}
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	if t, ok := c.last[key]; ok && now.Sub(t) < c.ttl {
		return false
	}
	c.last[key] = now
	if len(c.last) > 50_000 {
		c.gcLocked(now)
	}
	return true
}

func (c *CoolDown) gcLocked(now time.Time) {
	for k, t := range c.last {
		if now.Sub(t) > c.ttl*2 {
			delete(c.last, k)
		}
	}
}
