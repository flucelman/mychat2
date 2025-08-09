package controller

import (
	"backend/config"
	"backend/global"
	"backend/models"
	"backend/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
查询聊天记录表
通过user_id查询聊天记录
*/
func GetChatHistory(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	chatHistory := []models.ChatHistory{}

	// 在查询时直接按 updated_at 降序排列（最新的在前）
	if err := global.DB.Where("user_id = ?", userID).Order("updated_at desc").Find(&chatHistory).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type ChatHistoryResponse struct {
		ChatID    string    `json:"chat_id"`
		Title     string    `json:"title"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	chatHistoryResponse := []ChatHistoryResponse{}
	for _, chat := range chatHistory {
		chatHistoryResponse = append(chatHistoryResponse, ChatHistoryResponse{
			ChatID:    chat.ChatID,
			Title:     chat.Title,
			UpdatedAt: chat.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"data": chatHistoryResponse})
}

/*
查询消息表
通过chat_id查询消息

	返回：{
		message_id: string,
		role: string,
		content: string,
		model: string,
		created_at: time.Time,
	}
*/
func GetChatMessage(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	chatID := ctx.Param("chat_id")

	// 从查询参数获取分页信息，设置默认值
	page := 1
	pageSize := 20

	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}
	message := []models.Message{}
	if err := global.DB.Where("user_id = ?", userID).Where("chat_id = ?", chatID).
		Order("created_at ASC"). // 按创建时间升序排列（最早的在前）
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&message).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type MessageResponse struct {
		MessageID string    `json:"message_id"`
		Role      string    `json:"role"`
		Content   string    `json:"content"`
		Model     string    `json:"model"`
		CreatedAt time.Time `json:"created_at"`
	}
	messageResponse := []MessageResponse{}
	for _, message := range message {
		messageResponse = append(messageResponse, MessageResponse{
			MessageID: message.MessageID,
			Role:      message.Role,
			Content:   message.Content,
			Model:     message.Model,
			CreatedAt: message.CreatedAt,
		})
	}
	ctx.JSON(http.StatusOK, messageResponse)
}

// 查询模型列表
func GetModelList(ctx *gin.Context) {
	modelList := config.AIModels
	ctx.JSON(http.StatusOK, modelList)
}

/*
添加对话信息并通过SSE流式返回AI响应
如果是空聊天记录，则创建新的聊天记录
如果非空，则添加对话信息
*/
func AddChatMessage(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var input struct {
		ChatID         string           `json:"chat_id"`
		AIConfig       map[string]any   `json:"AI_config"`
		MessageHistory []map[string]any `json:"message_history"`
		FileUrl        []string         `json:"file_url"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果chat_id为空，则创建新的聊天记录
	if input.ChatID == "" {
		input.ChatID = uuid.New().String()
		// 1. 创建新的聊天记录
		chat := models.ChatHistory{
			ChatID: input.ChatID,
			UserID: userID,
			Title: func() string {
				content := input.MessageHistory[1]["content"].(string)
				if len(content) <= 30 {
					return content
				}
				return content[:30] + "..."
			}(),
		}
		if err := global.DB.Create(&chat).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	// 2. 添加用户信息到数据库
	fmt.Println("用户消息ID：", input.MessageHistory[len(input.MessageHistory)-1]["message_id"].(string))
	response := utils.SaveDB(input.MessageHistory[len(input.MessageHistory)-1]["message_id"].(string), userID, input.ChatID, "user", input.MessageHistory[len(input.MessageHistory)-1]["content"].(string), input.AIConfig["model"].(string))
	if response != "success" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": response})
		return
	}

	// 3.检查是否有file_url
	if len(input.FileUrl) > 0 {
		// 拿到url，进行解析
		// 如果是文件则进行markitdown返回string，如果是图片、视频、音频则返回url
		// 将string或者url保存到数据库
		// 调用不一样的AI接口
		fmt.Println("开始解析文件")
	}

	// 4. 设置SSE响应头
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Cache-Control")

	// 5. 发送初始响应（包含chat_id）
	assistantMessageID := uuid.New().String()
	ctx.SSEvent("start", gin.H{"chat_id": input.ChatID, "assistant_message_id": assistantMessageID})
	ctx.Writer.Flush()

	// 6. 参数验证和类型转换
	temperature, ok := input.AIConfig["temperature"].(float64)
	if !ok {
		ctx.SSEvent("error", gin.H{"error": "temperature must be a number"})
		return
	}

	topP, ok := input.AIConfig["top_p"].(float64)
	if !ok {
		ctx.SSEvent("error", gin.H{"error": "top_p must be a number"})
		return
	}

	frequencyPenalty, ok := input.AIConfig["frequency_penalty"].(float64)
	if !ok {
		ctx.SSEvent("error", gin.H{"error": "frequency_penalty must be a number"})
		return
	}

	maxTokens, ok := input.AIConfig["max_tokens"].(float64)
	if !ok {
		ctx.SSEvent("error", gin.H{"error": "max_tokens must be a number"})
		return
	}

	// 7. 创建通道和上下文
	answerCh := make(chan string)
	answerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 8. 启动AI流式响应
	modelKey := config.GetModelKey(input.AIConfig["model"].(string))
	go utils.AIStreamResponse(answerCtx, answerCh, modelKey, float32(temperature), int(maxTokens), float32(topP), float32(frequencyPenalty), input.MessageHistory)

	// 9. 处理流式响应并通过SSE发送
	fullResponse := ""
	for {
		select {
		case content, ok := <-answerCh:
			if !ok {
				// 通道关闭，AI响应结束
				// 保存完整的AI响应到数据库
				fmt.Println("助手消息ID：", assistantMessageID)
				saveResponse := utils.SaveDB(assistantMessageID, userID, input.ChatID, "assistant", fullResponse, input.AIConfig["model"].(string))

				if saveResponse != "success" {
					ctx.SSEvent("error", gin.H{"error": saveResponse})
				} else {
					ctx.SSEvent("end", gin.H{"message": "AI响应完成", "message_id": assistantMessageID})
				}
				return
			}
			// 发送内容片段
			ctx.SSEvent("content", gin.H{"content": content})
			fullResponse += content
			ctx.Writer.Flush()

		case <-ctx.Request.Context().Done():
			// 客户端断开连接
			fmt.Println("客户端断开连接，停止AI响应")
			if fullResponse == "" {
				fmt.Println("助手消息为空，不保存到数据库")
				cancel()
				return
			}
			fmt.Println("助手消息ID：", assistantMessageID)
			saveResponse := utils.SaveDB(assistantMessageID, userID, input.ChatID, "assistant", fullResponse, input.AIConfig["model"].(string))
			if saveResponse != "success" {
				ctx.SSEvent("error", gin.H{"error": saveResponse})
			} else {
				ctx.SSEvent("end", gin.H{"message": "AI响应完成", "message_id": assistantMessageID})
			}
			cancel()
			return
		}
	}
}

// 删除所有历史记录
func DeleteAllHistory(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if err := global.DB.Where("user_id = ?", userID).Delete(&models.ChatHistory{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := global.DB.Where("user_id = ?", userID).Delete(&models.Message{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "所有历史记录已删除"})
}

// 删除单个聊天记录
func DeleteSingleHistory(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var input struct {
		ChatID string `json:"chat_id"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := global.DB.Where("user_id = ?", userID).Where("chat_id = ?", input.ChatID).Delete(&models.ChatHistory{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := global.DB.Where("user_id = ?", userID).Where("chat_id = ?", input.ChatID).Delete(&models.Message{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "聊天记录已删除"})
}
