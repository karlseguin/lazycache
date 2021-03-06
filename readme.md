### A Lazy Cache storage for Go

`LazyCache` reloads all values on a given interval while lazilly loading new values on demand. Its simplicity is ideal for storing a small set of objects, where the cost of bulk-reloading is smaller than maintaining a more complex caching algorithm.

`LazyCache` requires two key pieces: a `Fetcher` and a `Loader`. The `Fetcher` is used to get individually items. The Loader is used to bulk reload all values.

For example:

```go
func fetch(id string) (interface{}, error) {
  account := new(Account)
  var name string
  row := db.QueryRow("selectname from accounts where id = ?", id)
  if err := row.Scan(&name); err != nil {
    // err will get returned by the call to Get which caused
    // this fetch
    return nil, err
  }
  return &Account{id, name}
}

func reload() (map[string]interface{}, error) {
  rows, err := db.Query("select id, name from accounts")
  if err != nil {
    //todo: log
    return nil, err
  }
  accounts := make(map[string]interface{})
  for rows.Next() {
    var id, name string
    rows.Scan(&id, &name)
    accounts[id] = &Account{id, name}
  }
  return accounts, nil
}

cache := lazycache.New(fetch, reload, time.Minute * 2)
```

When you `Get` from the cache, `fetch` will be invoked to load the item. Every 2 minutes, `reload` will be invoked to update the cache. This hybrid approach means that you can keep data in memory without having to wait until the next reload to discover new items.

See [ccache](https://github.com/karlseguin/ccache) for more typical and powerful LRU cache.

## MissGuard

`Fetcher` is called whenever the item isn't found in the cache. In many cases, you may wish to cache a miss for a short period. `MissGuard` is a wrapper around your own `Fetcher` which caches misses for a defined TTL:

```go
func yourFetcher(id string) (interface{}, error) {
  //hit the db
  if found == false {
    return nil, nil
  }
  return found, nil
}

// store up to 1000 misses, and each miss is cached for 5 seconds
mg := lazycache.NewMissGuard(yourFetcher, 1000, time.Second * 5)
cache := lazycache.New(mg.Fetch, reload, time.Minute * 2)
```

`MissGuard` only caches misses when your own fetcher returns a nil value AND a nil error.
