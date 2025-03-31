package memory

import (
	"github.com/duanchi/min/v2/abstract"
	"github.com/patrickmn/go-cache"
	"net/url"
	"time"
)

type MemoryCache struct {
	abstract.Bean
	instance *cache.Cache
}

func (this *MemoryCache) Init() {

}

func (this *MemoryCache) Instance(dsn *url.URL) {
	this.instance = cache.New(cache.NoExpiration, 1*time.Minute)
	// this.instance = cache2go.Cache(dsn.Hostname())
}

func (this *MemoryCache) Get(key string) (value interface{}) {
	value, _ = this.instance.Get(key)
	return
}

func (this *MemoryCache) Has(key string) bool {
	_, has := this.instance.Get(key)

	return has
}

func (this *MemoryCache) Set(key string, value interface{}) {
	this.instance.Set(key, value, cache.NoExpiration)
}

func (this *MemoryCache) SetWithTTL(key string, value interface{}, ttl int) {
	if ttl <= 0 {
		ttl = 0
	}

	this.instance.Set(key, &value, time.Duration(ttl)*time.Second)
}

func (this *MemoryCache) Del(key string) {
	this.instance.Delete(key)
}

func (this *MemoryCache) Flush() {
	this.instance.Flush()
}
