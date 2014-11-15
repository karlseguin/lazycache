package lazycache

import (
	"sync"
	"time"
)

type Fetcher func(key string) (interface{}, error)
type Loader func() (map[string]interface{}, error)

type LazyCache struct {
	sync.RWMutex
	fetcher Fetcher
	loader  Loader
	ttl     time.Duration
	items   map[string]interface{}
}

func New(fetcher Fetcher, loader Loader, ttl time.Duration) *LazyCache {
	cache := &LazyCache{
		ttl:     ttl,
		loader:  loader,
		fetcher: fetcher,
		items:   make(map[string]interface{}),
	}
	go cache.reloader()
	return cache
}

func (c *LazyCache) Get(key string) (interface{}, error) {
	c.RLock()
	item, exists := c.items[key]
	c.RUnlock()
	if exists == false {
		return c.fetch(key)
	}
	return item, nil
}

func (c *LazyCache) Set(key string, item interface{}) {
	c.Lock()
	defer c.Unlock()
	c.items[key] = item
}

func (c *LazyCache) Reload() {
	items, err := c.loader()
	if err != nil {
		return
	}
	c.Lock()
	c.items = items
	c.Unlock()
}

func (c *LazyCache) reloader() {
	for {
		time.Sleep(c.ttl)
		c.Reload()
	}
}

func (c *LazyCache) fetch(key string) (interface{}, error) {
	item, err := c.fetcher(key)
	if err != nil {
		return nil, err
	}
	c.Set(key, item)
	return item, nil
}
