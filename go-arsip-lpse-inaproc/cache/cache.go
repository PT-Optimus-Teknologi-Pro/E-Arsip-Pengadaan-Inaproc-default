package cache

import (
	"time"

	"github.com/gofiber/fiber/v2/log"
	cache "github.com/patrickmn/go-cache"
)

var internalCache *cache.Cache

func init() {
	log.Info("initial cache setup...")
	internalCache = cache.New(1*time.Hour, 10*time.Minute)
}

func Set(k string, x interface{}) {
	internalCache.SetDefault(k, x)
}

func Get(k string) (interface{}, bool) {
	return  internalCache.Get(k)
}

func Delete(k string) {
	internalCache.Delete(k)
}

func Flush() {
	internalCache.Flush()
}
