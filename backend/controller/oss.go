package controller

import (
	"net/http"
	"sync"

	"backend/utils"

	"github.com/gin-gonic/gin"
)

func UploadFile(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	files := form.File["files"]

	// 并发上传，限制同时进行的上传数为3个
	fileURLs := make([]string, len(files))
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup
	var firstErr error
	var setErrOnce sync.Once

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
				setErrOnce.Do(func() { firstErr = err })
				return
			}
			fileURLs[i] = url
		}()
	}
	wg.Wait()

	if firstErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": firstErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "File uploaded successfully", "file_urls": fileURLs})
}
