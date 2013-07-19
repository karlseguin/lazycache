package lazycache

import (
  "sync"
  "time"
)

type Fetcher func(id string) (interface{})

type Item struct {
  object interface{}
  expires time.Time
}

type LazyCache struct {
  fetcher Fetcher
  ttl time.Duration
  lock *sync.RWMutex
  items map[string]*Item
}

func New(fetcher Fetcher, ttl time.Duration, size int) *LazyCache {
  return &LazyCache{
    ttl: ttl,
    fetcher: fetcher,
    lock: new(sync.RWMutex),
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
  return item.object, true
}

func (cache *LazyCache) Fetch(id string, current *Item) (interface{}, bool) {
  object := cache.fetcher(id)
  if object != nil || current == nil {
    cache.lock.Lock()
    if current != nil {
      current.expires = time.Now().Add(cache.ttl)
      current.object = object
    } else {
      cache.items[id] = &Item{ expires: time.Now().Add(cache.ttl), object: object}
    }
    cache.lock.Unlock()
  }
  return object, object != nil
}

