package lazycache

import (
  "time"
  "errors"
  "testing"
)

func TestFetchesWhenNotCached(t *testing.T) {
  count := 0
  cache := New(countFetcher(&count), time.Second, 1)
  cache.Get("Hi")
  if count != 1 {
    t.Errorf("expected %+v to equal 1", count)
  }
}

func TestDoesNotFetchWhenCached(t *testing.T) {
  count := 0
  cache := New(countFetcher(&count), time.Second, 1)
  cache.Get("Hi")
  cache.Get("Hi")
  if count != 1 {
    t.Errorf("expected %+v to equal 1", count)
  }
}

func TestReturnsCachedAndFetchesLazilyAfterTtl(t *testing.T) {
  count := 0
  cache := New(countFetcher(&count), time.Microsecond, 1)

  cache.Get("Hi")

  // Second get, returns old value (1) and fetches on the background.
  time.Sleep(5 * time.Microsecond)
  v2, _ := cache.Get("Hi")

  if v2.(int) != 1 {
    t.Errorf("expected %+v to equal 1", v2.(int))
  }

  time.Sleep(2 * time.Microsecond)

  if count != 2 {
    t.Errorf("expected %+v to equal 2", count)
  }

  v3, _ := cache.Get("Hi")

  if v3.(int) != 2 {
    t.Errorf("expected %+v to equal 2", v3.(int))
  }
}

func TestDoesNotFetchErrorsUntilExpire(t *testing.T) {
  count := 0
  cache := New(nilFetcher(&count), time.Second, 1)
  cache.Get("Hi")
  cache.Get("Hi")
  if count != 1 {
    t.Errorf("expected %+v to equal 1", count)
  }
}

func TestFetchesErrorsSynchronouslyAfterExpire(t *testing.T) {
  count := 0
  cache := New(slowNilFetcher(&count), time.Microsecond, 1)

  cache.Get("Hi")
  time.Sleep(2 * time.Microsecond)
  cache.Get("Hi")

  if count != 2 {
    t.Errorf("expected %+v to equal 2", count)
  }

  time.Sleep(2 * time.Microsecond)
  cache.Get("Hi")

  if count != 3 {
    t.Errorf("expected %+v to equal 3", count)
  }
}

func TestFetchingNilErasesExistingValue(t *testing.T) {
  count := 0
  cache := New(nilFetcher(&count), time.Microsecond, 1)
  cache.items["Hi"] = &Item{object: 99, expires: time.Now(),}

  v1, _ := cache.Get("Hi")

  if v1.(int) != 99 {
    t.Errorf("expected %+v to equal 99", v1.(int))
  }
  time.Sleep(2 * time.Microsecond)

  v2, ok := cache.Get("Hi")

  if ok != false && v2 != nil {
    t.Errorf("expected value to be removed from the cache")
  }
}

func TestErrorOnFetchKeepsOldValue(t *testing.T) {
  count := 0
  cache := New(errorFetcher(&count), time.Microsecond, 1)
  cache.items["paul"] = &Item{object: 99, expires: time.Now().Add(-time.Hour),}

  v1, _ := cache.Get("paul")

  if v1.(int) != 99 {
    t.Errorf("expected %+v to equal 99", v1.(int))
  }
}

func countFetcher(count *int) Fetcher {
  return func (id string) (interface{}, error) {
    *count += 1
    return *count, nil
  }
}

func nilFetcher(count *int) Fetcher {
  return func (id string) (interface{}, error) {
    *count += 1
    return nil, nil
  }
}

func slowNilFetcher(count *int) Fetcher {
  return func (id string) (interface{}, error) {
    time.Sleep(10 * time.Microsecond)
    *count += 1
    return nil, nil
  }
}

func errorFetcher(count *int) Fetcher {
  return func (id string) (interface{}, error) {
    return nil, errors.New("oops")
  }
}