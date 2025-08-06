package controller

import (
	"backend/global"
	"backend/models"
	"backend/utils"
	"context"
	"fmt"
	"net/http"
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
	if err := global.DB.Where("user_id = ?", userID).Find(&chatHistory).Error; err != nil {
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
	ctx.JSON(http.StatusOK, chatHistoryResponse)
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
	type MessageRequest struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	messageRequest := MessageRequest{}
	if err := ctx.ShouldBindJSON(&messageRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message := []models.Message{}
	if err := global.DB.Where("user_id = ?", userID).Where("chat_id = ?", chatID).Offset((messageRequest.Page - 1) * messageRequest.PageSize).Limit(messageRequest.PageSize).Find(&message).Error; err != nil {
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

/*
添加对话信息并通过SSE流式返回AI响应
如果是空聊天记录，则创建新的聊天记录
如果非空，则添加对话信息
*/
func AddUserMessage(ctx *gin.Context) {
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
				if len(input.MessageHistory[0]["content"].(string)) <= 10 {
					return input.MessageHistory[0]["content"].(string)
				}
				return input.MessageHistory[0]["content"].(string)[:10] + "..."
			}(),
		}
		if err := global.DB.Create(&chat).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// 2. 添加用户信息到数据库
	response := utils.SaveDB(userID, input.ChatID, "user", input.MessageHistory[len(input.MessageHistory)-1]["content"].(string), input.AIConfig["model"].(string))
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
	ctx.SSEvent("start", gin.H{"chat_id": input.ChatID})
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
	answerCh := make(chan string, 50)
	answerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 8. 启动AI流式响应
	go utils.AIStreamResponse(answerCtx, answerCh, input.AIConfig["model"].(string), float32(temperature), int(maxTokens), float32(topP), float32(frequencyPenalty), input.MessageHistory)

	// 9. 处理流式响应并通过SSE发送
	fullResponse := ""
	for {
		select {
		case content, ok := <-answerCh:
			if !ok {
				// 通道关闭，AI响应结束
				// 保存完整的AI响应到数据库
				saveResponse := utils.SaveDB(userID, input.ChatID, "assistant", fullResponse, input.AIConfig["model"].(string))
				if saveResponse != "success" {
					ctx.SSEvent("error", gin.H{"error": saveResponse})
				} else {
					ctx.SSEvent("end", gin.H{"message": "AI响应完成"})
				}
				return
			}
			// 发送内容片段
			fullResponse += content
			ctx.SSEvent("content", gin.H{"content": content})
			ctx.Writer.Flush()

		case <-ctx.Request.Context().Done():
			// 客户端断开连接
			fmt.Println("客户端断开连接，停止AI响应")
			cancel()
			return
		}
	}
}

func AddAssistantMessage(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var input struct {
		ChatID           string `json:"chat_id"`
		AssistantMessage string `json:"assistant_message"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response := utils.SaveDB(userID, input.ChatID, "assistant", input.AssistantMessage, "assistant")
	if response != "success" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": response})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"chat_id": input.ChatID}})
}
