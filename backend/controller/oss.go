package controller

import (
	"backend/global"
	"backend/models"
	"backend/utils"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	files := form.File["files"]

	// 并发上传，限制同时进行的上传数为10个
	successInfos := make([]gin.H, len(files))
	failedInfos := make([]gin.H, len(files))
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup

	for i, f := range files {
		i, f := i, f // 捕获当前循环变量
		wg.Add(1)
		go func() {
			defer wg.Done()

			// 限制并发
			sem <- struct{}{}
			defer func() { <-sem }()

			url, err := utils.UploadFile(f, "chat-files", userID)
			if err != nil {
				suffix := ""
				if dot := strings.LastIndex(f.Filename, "."); dot != -1 {
					suffix = f.Filename[dot:]
				}
				failedInfos[i] = gin.H{
					"name":   f.Filename,
					"suffix": suffix,
					"size":   f.Size,
					"error":  err.Error(),
				}
				return
			}
			fileID := uuid.New().String()
			// url 写到数据库
			global.DB.Create(&models.File{
				FileID:   fileID,
				UserID:   userID,
				FileName: f.Filename,
				FileSize: f.Size,
				FileURL:  url,
			})
			suffix := ""
			if dot := strings.LastIndex(f.Filename, "."); dot != -1 {
				suffix = f.Filename[dot:]
			}
			successInfos[i] = gin.H{
				"name":    f.Filename,
				"suffix":  suffix,
				"size":    f.Size,
				"file_id": fileID,
			}
		}()
	}
	wg.Wait()

	// 过滤掉未填充的项
	filteredSuccess := make([]gin.H, 0, len(files))
	for _, info := range successInfos {
		if info != nil {
			filteredSuccess = append(filteredSuccess, info)
		}
	}
	filteredFailed := make([]gin.H, 0, len(files))
	for _, info := range failedInfos {
		if info != nil {
			filteredFailed = append(filteredFailed, info)
		}
	}

	if len(filteredSuccess) == 0 && len(filteredFailed) > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "所有文件上传失败",
			"data": gin.H{
				"success": filteredSuccess,
				"failed":  filteredFailed,
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文件上传完成",
		"data": gin.H{
			"success": filteredSuccess,
			"failed":  filteredFailed,
		},
	})
}
