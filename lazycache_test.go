package lazycache

import (
	// "errors"
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
	cache := New(CloneFetcher, nil, time.Minute)
	Expect(cache.Get("spice")).To.Equal("fetch:spice", nil)
}

func CloneFetcher(key string) (interface{}, error) {
	return "fetch:" + key, nil
}
