package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var redisInstance *RedisPool

type RedisPool struct {
	client *redis.Client
}

func (self *RedisPool) init() {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", "", ""),
		Password: "",
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	self.client = client
}
func NewRedis() *RedisPool {
	if redisInstance == nil {
		redisInstance = &RedisPool{}
		redisInstance.init()
	}
	return redisInstance
}

func (self *RedisPool) RHSet(key string, mapping map[string]interface{}, expire time.Duration) string {
	self.client.HSet(context.Background(), key, mapping)
	if expire > 0 {
		self.client.Expire(context.Background(), key, expire)
	}
	return key
}
func (self *RedisPool) RHGetAll(keys []string) []map[string]string {
	pipeline := self.client.Pipeline()
	for _, key := range keys {
		pipeline.HGetAll(context.Background(), key)
	}
	results, _ := pipeline.Exec(context.Background())

	nonEmptyResults := make([]map[string]string, 0)
	for _, result := range results {
		if result != nil {
			//nonEmptyResults = append(nonEmptyResults, result.(map[string]string))
		}
	}
	return nonEmptyResults
}
