package lazycache

import (
	"sync"
	"time"
)

type MissGuard struct {
	sync.RWMutex
	fetcher Fetcher
	ttl     time.Duration
	lookup  map[string]time.Time
}

func NewMissGuard(fetcher Fetcher, size int, ttl time.Duration) *MissGuard {
	return &MissGuard{
		ttl:     ttl,
		fetcher: fetcher,
		lookup:  make(map[string]time.Time),
	}
}

func (m *MissGuard) Fetch(key string) (interface{}, error) {
	now := time.Now()
	m.RLock()
	e, ok := m.lookup[key]
	m.RUnlock()

	if ok && e.After(now) {
		return nil, nil
	}

	// either it isn't in our lookup or it's expired
	item, err := m.fetcher(key)
	if item != nil || err != nil {
		return item, err
	}

	// the real fetcher returned nil, store it
	m.Lock()
	m.lookup[key] = now.Add(m.ttl)
	m.Unlock()
	return nil, nil
}
