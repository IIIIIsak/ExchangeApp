package config

import (
	"context"
	"exchangeapp/global"
	"github.com/redis/go-redis/v9"
	"log"
)

func InitRedis() {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     AppConig.Redis.Addr,
		Password: AppConig.Redis.Password,
		DB:       AppConig.Redis.DB,
	})

	ctx := context.Background()
	// 检查是否成功连接
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis, got error: %v", err)
	}

	global.RedisDB = RedisClient

}
