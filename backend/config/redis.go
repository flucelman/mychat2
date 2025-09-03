package config

import (
	"backend/global"
	"os"

	"github.com/go-redis/redis"
)

func InitRedis() {
	global.RedisClient = redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Host + ":" + AppConfig.Redis.Port,
		DB:       AppConfig.Redis.DB,
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	_, err := global.RedisClient.Ping().Result()
	if err != nil {
		panic("Failed to connect to Redis:" + err.Error())
	}
}
