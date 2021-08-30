package util

import (
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"os"
	"reflect"
	"time"
	"unsafe"
)

func GetType (i interface{}) reflect.Type {
	return reflect.TypeOf(i).Elem()
}

func Getenv (key string, defaults string) string {
	result := os.Getenv(key)

	if result == "" {
		return defaults
	} else {
		return result
	}
}

func GenerateUUID() uuid.UUID {
	return uuid.NewV4()
}




func RandomString(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxMask := int64(1<<6 - 1) // All 1-bits, as many as 6
	b := make([]byte, n)
	src := rand.NewSource(time.Now().UnixNano())

	// A src.Int63() generates 63 random bits, enough for 10 characters!
	for i, cache, remain := n-1, src.Int63(), 10; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), 10
		}
		b[i] = letterBytes[int(cache & letterIdxMask) % len(letterBytes)]
		i--
		cache >>= 6
		remain--
	}
	return  *(*string)(unsafe.Pointer(&b))
}