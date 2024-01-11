package redis

import (
    "context"
    "github.com/redis/go-redis/v9"
    "sync"
)

var RedisClient *redis.Client
var redisOnce sync.Once

func LoadRedis() *redis.Client {
    redisOnce.Do(func() {
        RedisClient = redis.NewClient(&redis.Options{
            Addr:            `AppConfig.GetString("redis.addr")`,
            Password:        `AppConfig.GetString("redis.addr")`,
            DB:              0,
            MaxRetries:      0,
            MinRetryBackoff: 0,
            MaxRetryBackoff: 0,
            DialTimeout:     0,
            ReadTimeout:     0,
            WriteTimeout:    0,
            PoolSize:        500,
            MinIdleConns:    50,
            PoolTimeout:     0,
            TLSConfig:       nil,
        })
        ctx := context.Background()
        if err := RedisClient.Ping(ctx).Err(); err != nil {
            panic(err.Error())
        }
    })
    return RedisClient
}
