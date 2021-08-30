package _interface

import "net/url"

type CacheInterface interface {
	Instance(dsn *url.URL)
	Has(key string) bool
	Get(key string) (value interface{})
	Set(key string, value interface{})
	SetWithTTL(key string, value interface{}, ttl int)
	Flush()
	Del(key string)
}
