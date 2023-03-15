package app

import (
	"context"
	"fmt"
	"github.com/frchandra/chatin/config"
	"github.com/redis/go-redis/v9"
)

func NewCache(appConfig *config.AppConfig) *redis.Client {
	cache := redis.NewClient(&redis.Options{
		Password: appConfig.RedisPassword,
		Addr:     appConfig.RedisHost + ":" + appConfig.RedisPort,
	})
	var ctx = context.Background()
	_, err := cache.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	//uncomment this during development phase
	cache.FlushAll(ctx)
	fmt.Println("redis connected")
	return cache

}
