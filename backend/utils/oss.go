package utils

import (
	"context"
	"errors"
	"mime/multipart"
	"os"

	"backend/global"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

func UploadFile(file *multipart.FileHeader, folder string, user_id string) (string, error) {
	// 基础校验，便于定位问题
	bucket := os.Getenv("OSS_BUCKET")
	base_url := os.Getenv("OSS_BASE_URL")

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	oss_key := folder + "/" + user_id + "/" + file.Filename

	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(bucket),  // 存储空间名称
		Key:    oss.Ptr(oss_key), // 对象名称
		Body:   src,              // 传入文件内容
	}
	_, err = global.OSSClient.PutObject(context.Background(), request)
	if err != nil {
		return "", errors.New(err.Error())
	}
	file_url := base_url + "/" + oss_key
	return file_url, nil
}
