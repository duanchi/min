package cache

import (
	"github.com/duanchi/min/cache/memory"
	"github.com/duanchi/min/cache/redis"
	"github.com/duanchi/min/context"
	_interface "github.com/duanchi/min/interface"
	"net/url"
	"reflect"
)

var engine _interface.CacheInterface

var CacheEngines map[string]reflect.Value

func Init() {
	CacheEngines = map[string]reflect.Value{}
	CacheEngines["memory"] = reflect.ValueOf(&memory.MemoryCache{})
	CacheEngines["redis"] = reflect.ValueOf(&redis.RedisCache{})

	dsn := context.GetApplicationContext().GetConfig("Cache.Dsn").(string)
	dsnUrl, _ := url.Parse(dsn)

	if _, ok := CacheEngines[dsnUrl.Scheme]; !ok {
		dsnUrl, _ = url.Parse("memory://min")
	}

	engine = CacheEngines[dsnUrl.Scheme].Interface().(_interface.CacheInterface)

	engine.Instance(dsnUrl)
}

func Has(key string) bool {
	return engine.Has(key)
}

func Get(key string) (value interface{}) {
	return engine.Get(key)
}

func Set(key string, value interface{}) {
	engine.Set(key, value)
}

func SetWithTTL(key string, value interface{}, ttl int) {
	engine.SetWithTTL(key, value, ttl)
}

func Del(key string) {
	engine.Del(key)
}

func Flush() {
	engine.Flush()
}
