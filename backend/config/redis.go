package config

import (
	"backend/global"

	"github.com/go-redis/redis"
)

func InitRedis() {
	global.RedisClient = redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Host + ":" + AppConfig.Redis.Port,
		DB:       AppConfig.Redis.DB,
		Password: AppConfig.Redis.Password,
	})

	_, err := global.RedisClient.Ping().Result()
	if err != nil {
		panic("Failed to connect to Redis:" + err.Error())
	}
}
