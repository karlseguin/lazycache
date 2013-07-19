package lazycache

import (
  "time"
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

func TestReturnsOldValuesWhenGettingErrors(t *testing.T) {
  count := 0
  cache := New(nilFetcher(&count), time.Microsecond, 1)
  cache.items["Hi"] = &Item{object: 99, expires: time.Now(),}

  v1, _ := cache.Get("Hi")

  if v1.(int) != 99 {
    t.Errorf("expected %+v to equal 99", v1.(int))
  }
  time.Sleep(2 * time.Microsecond)

  v2, _ := cache.Get("Hi")

  if v2.(int) != 99 {
    t.Errorf("expected %+v to equal 99", v2.(int))
  }
}

func TestLol(t *testing.T){
  fetcher := func (id string) (interface{}) {
    return id == "foo"
  }

  cache := New(fetcher, 60 * time.Second, 256) // Prealocates 256 items. Items are expired after 60s.

  foo_value, foo_found := cache.Get("foo")
  print(foo_found)
  print(foo_value.(bool))
  bar_value, bar_found := cache.Get("Bar")
  print(bar_found)
  print(bar_value.(bool))
  t.Errorf("expected")
}

func countFetcher(count *int) Fetcher {
  return func (id string) (interface{}) {
    *count += 1
    return *count
  }
}
func nilFetcher(count *int) Fetcher {
  return func (id string) (interface{}) {
    *count += 1
    return nil
  }
}
func slowNilFetcher(count *int) Fetcher {
  return func (id string) (interface{}) {
    time.Sleep(10 * time.Microsecond)
    *count += 1
    return nil
  }
}
