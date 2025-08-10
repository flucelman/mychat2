package global

import (
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"gorm.io/gorm"

	"github.com/go-redis/redis"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
	OSSClient   *oss.Client
)
