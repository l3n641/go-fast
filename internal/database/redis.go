package database

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"sync"
)

type RedisDatabase struct {
	client *redis.Client
	once   sync.Once
}

func (r *RedisDatabase) GetInstance() *redis.Client {
	r.once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:         viper.GetString("redis.addr"),
			Password:     viper.GetString("redis.password"),
			DB:           viper.GetInt("redis.DB"),
			PoolSize:     viper.GetInt("redis.poolSize"),
			MinIdleConns: viper.GetInt("redis.minIdleConns"),
		})
		r.client = client

	})

	return r.client
}

func NewRedisClient() *RedisDatabase {
	return &RedisDatabase{}
}
