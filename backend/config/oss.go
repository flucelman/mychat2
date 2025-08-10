package config

import (
	"os"

	"backend/global"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

func InitOSS() {

	region := os.Getenv("OSS_REGION")

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	// 创建OSS客户端
	global.OSSClient = oss.NewClient(cfg)
}
