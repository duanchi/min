package redis

import (
	"context"
	"fmt"
	"github.com/duanchi/min/abstract"
	"github.com/go-redis/redis/v8"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type RedisCache struct {
	abstract.Bean
	instance *redis.Client
	ctx      context.Context
}

func (this *RedisCache) Init() {
	this.ctx = context.Background()
}

func (this *RedisCache) Ctx() context.Context {
	if this.ctx == nil {
		this.ctx = context.Background()
	}

	return this.ctx
}

func (this *RedisCache) Instance(dsn *url.URL) {
	password, _ := dsn.User.Password()
	pathString := strings.Trim(dsn.Path, "/")
	if pathString == "" || pathString == "/" {
		pathString = "0"
	}
	path, _ := strconv.Atoi(pathString)
	this.instance = redis.NewClient(&redis.Options{
		Addr:     dsn.Host,
		Password: password, // no password set
		DB:       path,     // use default DB
	})
	fmt.Printf("Redis %s connected at DB %d!\r\n", dsn.Host, path)
}

func (this *RedisCache) Get(key string) (value interface{}) {
	value, _ = this.instance.Get(this.Ctx(), key).Result()
	return
}

func (this *RedisCache) Has(key string) bool {

	has, _ := this.instance.Exists(this.Ctx(), key).Result()

	if has > 0 {
		return true
	}
	return false
}

func (this *RedisCache) Set(key string, value interface{}) {
	this.instance.Set(this.Ctx(), key, value, 0).Result()
}

func (this *RedisCache) SetWithTTL(key string, value interface{}, ttl int) {
	if ttl <= 0 {
		ttl = 0
	}

	this.instance.Set(this.Ctx(), key, value, time.Duration(ttl)*time.Second).Result()
}

func (this *RedisCache) Del(key string) {
	this.instance.Del(this.Ctx(), key).Result()
}

func (this *RedisCache) Flush() {
	this.instance.FlushDB(this.Ctx()).Result()
}
