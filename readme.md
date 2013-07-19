### A Lazy Cache storage for Go

Lazy cache is a thread-safe key-value storage that fetches the new values asynchronously when possible.

* When the cache doesn't have any value for a given key, it will fetch synchronously and return the new value.
* When the cache have a value for a key and it's still valid, it will return it.
* When the cache have a value for a key and it's outdated, it will return it and fetch the new one on the background.

It also handles errors. When the fetcher function returns `nil`, it will save it on the cache, but will never return a stale version.

### Example

    fetcher := func (id string) (interface{}) {
      return id == "foo"
    }

    cache := New(fetcher, 60 * time.Second, 256) // Prealocates 256 items. Items are expired after 60s.

    foo_value, foo_found := cache.Get('foo')
    bar_value, bar_found := cache.Get('bar')


### Installation
Install using the "go get" command:

    go get github.com/viki-org/lazycache
