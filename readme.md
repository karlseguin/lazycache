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
    return nil, err
  }
  return &Account{id, name}
}

func reload() []interface{} {
  rows, err := db.Query("select id, name from accounts")
  if err != nil {
    //todo: log
    return nil
  }
  accounts := make([]Account, 0, 10)
  for rows.Next() {
    var id int
    var name string
    rows.Scan(&id, &name)
    accounts = append(accounts, &Account{id, name})
  }
  return accounts
}

cache := lazycache.New(fetch, reload, time.Minute * 2)
```

When you `Get` from the cache, `fetch` will be invoked to load the item. Every 2 minutes, `reload` will be invoked to update the cache. This hybrid approach means that you can keep data in memory without having to wait until the next reload to discover new items.

See [ccache](https://github.com/karlseguin/ccache) for more typical and powerful LRU cache.
