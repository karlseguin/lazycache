package lazycache

import (
	"errors"
	. "github.com/karlseguin/expect"
	"testing"
	"time"
)

type LazyCacheTests struct{}

func Test_LazyCache(t *testing.T) {
	Expectify(new(LazyCacheTests), t)
}

func (_ LazyCacheTests) GetsFromTheCache() {
	cache := New(nil, nil, time.Minute)
	cache.Set("leto", "ghanima")
	Expect(cache.Get("leto")).To.Equal("ghanima", nil)
}

func (_ LazyCacheTests) GetDoesAFetchOnMiss() {
	cache := New(cloneFetcher, nil, time.Minute)
	Expect(cache.Get("spice")).To.Equal("fetch:spice", nil)
}

func (_ LazyCacheTests) ReloadsWithNewItems() {
	cache := New(nil, loader(nil, "a", 1, "b", 2), time.Minute)
	cache.Reload()
	Expect(cache.Get("a")).To.Equal(1, nil)
	Expect(cache.Get("b")).To.Equal(2, nil)
	Expect(len(cache.items)).To.Equal(2)
}

func (_ LazyCacheTests) ReloadsWithEmptyItems() {
	cache := New(nil, loader(nil), time.Minute)
	cache.Set("a", 32)
	cache.Reload()
	Expect(len(cache.items)).To.Equal(0)
}

func (_ LazyCacheTests) ReloadSkipsOnError() {
	cache := New(nil, loader(errors.New("abc")), time.Minute)
	cache.Set("a", 32)
	cache.Reload()
	Expect(cache.Get("a")).To.Equal(32, nil)
	Expect(len(cache.items)).To.Equal(1)
}

func cloneFetcher(key string) (interface{}, error) {
	return "fetch:" + key, nil
}

func loader(err error, values ...interface{}) Loader {
	return func() (map[string]interface{}, error) {
		if err != nil {
			return nil, err
		}
		m := make(map[string]interface{}, len(values))
		for i := 0; i < len(values); i += 2 {
			m[values[i].(string)] = values[i+1]
		}
		return m, nil
	}
}
