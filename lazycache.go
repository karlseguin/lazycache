package lazycache
// A generic cache that favors returning stale data
// than blocking a caller

import (
  "sync"
  "time"
)

type Fetcher func(id string) (interface{}, error)

type Item struct {
  object interface{}
  expires time.Time
}

type LazyCache struct {
  fetcher Fetcher
  ttl time.Duration
  lock sync.RWMutex
  items map[string]*Item
}

func New(fetcher Fetcher, ttl time.Duration, size int) *LazyCache {
  return &LazyCache{
    ttl: ttl,
    fetcher: fetcher,
    items: make(map[string]*Item, size),
  }
}

func (cache *LazyCache) Get(id string) (interface{}, bool) {
  cache.lock.RLock()
  item, exists := cache.items[id]
  cache.lock.RUnlock()
  if exists == false {
    return cache.Fetch(id, item)
  }
  if time.Now().After(item.expires) {
    if item.object == nil {
      return cache.Fetch(id, item)
    } else {
      go cache.Fetch(id, item)
    }
  }
  return item.object, item.object != nil
}

func (cache *LazyCache) Set(id string, object interface{}) {
  cache.lock.Lock()
  defer cache.lock.Unlock()
  cache.items[id] = &Item{expires: time.Now().Add(cache.ttl), object: object}
}

func (cache *LazyCache) Fetch(id string, current *Item) (interface{}, bool) {
  object, err := cache.fetcher(id)
  if err != nil { return nil, false }
  
  if current == nil {
    cache.Set(id, object)
    return object, object != nil
  }
  cache.lock.Lock()
  defer cache.lock.Unlock()
  current.expires = time.Now().Add(cache.ttl)
  current.object = object
  return object, object != nil
}
