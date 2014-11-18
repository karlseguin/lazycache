package lazycache

import (
	"math/rand"
	"sync"
	"time"
)

type MissGuard struct {
	sync.RWMutex
	fetcher Fetcher
	size    int
	ttl     time.Duration
	slots   []time.Time
	lookup  map[string]int
}

func NewMissGuard(fetcher Fetcher, size int, ttl time.Duration) *MissGuard {
	return &MissGuard{
		ttl:     ttl,
		size:    size,
		fetcher: fetcher,
		slots:   make([]time.Time, size),
		lookup:  make(map[string]int, size),
	}
}

func (m *MissGuard) Fetch(key string) (interface{}, error) {
	var e time.Time
	now := time.Now()
	m.RLock()
	index, ok := m.lookup[key]
	if ok {
		e = m.slots[index]
	}
	m.RUnlock()

	if e.After(now) {
		return nil, nil
	}

	// either it isn't in our lookup or it's expired
	item, err := m.fetcher(key)
	if item != nil || err != nil {
		return item, err
	}

	// the real fetcher returned nil, store it
	index = rand.Intn(m.size)
	m.Lock()
	m.lookup[key] = index
	m.slots[index] = now.Add(m.ttl)
	m.Unlock()
	return nil, nil
}
