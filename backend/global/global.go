package global

import (
	"gorm.io/gorm"

	"github.com/go-redis/redis"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)
