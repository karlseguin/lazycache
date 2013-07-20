### A Lazy Cache storage for Go
Lazy cache is a key-value storage that favors returning stale data rather than blocking a caller. A cached item will always be returned, no matter how stale it is. However, expired items will be reloaded in a separate goroutine.

The only time the cache will block is when the key is unknown.

Should fetching an item return an error, an existing value will remain, even if stale.

### Example

    fetcher := func (id string) (interface{}, error) {
      return id == "foo", nil
    }

    // Prealocates 256 items. Items are expired after 2 minutes
    cache := New(fetcher, time.Minute * 2, 256) 

    foo_value, foo_found := cache.Get('foo')
    bar_value, bar_found := cache.Get('bar')


### Installation
Install using the "go get" command:

    go get github.com/viki-org/lazycache
