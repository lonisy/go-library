package go_cache

import (
    cache "github.com/patrickmn/go-cache"
    "time"
)

var Cache *cache.Cache

func init() {
    Cache = cache.New(5*time.Minute, 10*time.Minute)
}
