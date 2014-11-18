package lazycache

import (
	"fmt"
	. "github.com/karlseguin/expect"
	"testing"
	"time"
)

type MissGuardTests struct{}

func Test_MissGuard(t *testing.T) {
	Expectify(new(MissGuardTests), t)
}

func (_ *MissGuardTests) UsesFetcherIfNotGuarded() {
	m := NewMissGuard(cloneFetcher, 100, time.Minute)
	Expect(m.Fetch("spice")).To.Equal("fetch:spice", nil)
}

func (_ *MissGuardTests) UsesFetcherIfGuardExpired() {
	cf := NewCountFetcher("dune", "spice")
	m := NewMissGuard(cf.Fetch, 100, -time.Minute)
	Expect(m.Fetch("dune")).To.Equal("spice", nil)
	Expect(m.Fetch("dune")).To.Equal("spice", nil)
	Expect(cf.counts["dune"]).To.Equal(2)
}

func (_ *MissGuardTests) GuardsTheMiss() {
	cf := NewCountFetcher()
	m := NewMissGuard(cf.Fetch, 5, time.Minute)
	Expect(m.Fetch("dune")).To.Equal(nil, nil)
	Expect(m.Fetch("dune")).To.Equal(nil, nil)
	Expect(cf.counts["dune"]).To.Equal(1)
	fmt.Println(m.slots)
	fmt.Println(m.lookup)
}

type CountFetcher struct {
	values map[string]interface{}
	counts map[string]int
}

func NewCountFetcher(values ...interface{}) *CountFetcher {
	mapped := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		mapped[values[i].(string)] = values[i+1]
	}

	return &CountFetcher{
		values: mapped,
		counts: make(map[string]int),
	}
}

func (c *CountFetcher) Fetch(key string) (interface{}, error) {
	c.counts[key]++
	return c.values[key], nil
}
