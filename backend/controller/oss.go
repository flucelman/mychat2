package controller

import (
	"backend/global"
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 根据文件后缀获取文件类型
func getFileType(filename string) string {
	suffix := strings.ToLower(filename[strings.LastIndex(filename, "."):])
	switch suffix {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return "image"
	case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv":
		return "video"
	case ".mp3", ".wav", ".aac", ".m4a", ".ogg", ".flac":
		return "audio"
	case ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".md", ".html":
		return "file"
	default:
		return "file"
	}
}

// 上传文件
func UploadFile(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	chatID := form.Value["chat_id"][0]
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
			fileType := getFileType(f.Filename)
			fileID := uuid.New().String()
			// url 写到数据库
			global.DB.Create(&models.File{
				FileID:      fileID,
				UserID:      userID,
				ChatID:      chatID,
				FileName:    f.Filename,
				FileSize:    f.Size,
				FileURL:     url,
				FileType:    fileType,
				FileContent: "",
			})
			successInfos[i] = gin.H{
				"id":           fileID,
				"role":         "file",
				"file_name":    f.Filename,
				"file_size":    f.Size,
				"file_url":     url,
				"file_type":    fileType,
				"file_content": "",
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

// 删除文件（根据file_id）
func DeleteFile(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var input struct {
		FileID  string `json:"file_id"`
		FileURL string `json:"file_url"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 根据file_id删除数据库中的文件
	go func() {
		if err := global.DB.Delete(&models.File{}, "file_id = ? AND user_id = ?", input.FileID, userID).Error; err != nil {
			fmt.Println("删除数据库失败", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "文件删除失败",
			})
			return
		}
	}()

	// 根据file_url删除oss中的文件
	go func() {
		if err := utils.DeleteFile(input.FileURL); err != nil {
			fmt.Println("删除oss失败", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "文件删除失败",
			})
			return
		}
	}()

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文件删除成功",
	})
}
