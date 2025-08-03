package controller

import (
	"backend/global"
	"backend/models"
	"backend/utils"
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
添加对话信息
如果是空聊天记录，则创建新的聊天记录
如果非空，则添加对话信息
*/
func AddUserMessage(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var input struct {
		ChatID      string         `json:"chat_id"`
		AIConfig    map[string]any `json:"AI_config"`
		UserMessage string         `json:"user_message"`
		FileUrl     []string       `json:"file_url"`
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
				if len(input.UserMessage) <= 10 {
					return input.UserMessage
				}
				return input.UserMessage[:10] + "..."
			}(),
		}
		if err := global.DB.Create(&chat).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	// 2. 添加用户信息到数据库
	response := utils.SaveDB(userID, input.ChatID, "user", input.UserMessage, input.AIConfig["model"].(string))
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

	// 4. 调用普通对话AI接口
	fmt.Println("开始调用 AI 接口...")
	temperature, ok := input.AIConfig["temperature"].(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "temperature must be a number"})
		return
	}

	topP, ok := input.AIConfig["top_p"].(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "top_p must be a number"})
		return
	}

	frequencyPenalty, ok := input.AIConfig["frequency_penalty"].(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "frequency_penalty must be a number"})
		return
	}

	maxTokens, ok := input.AIConfig["max_tokens"].(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "max_tokens must be a number"})
		return
	}

	response = utils.AIResponse(
		input.AIConfig["model"].(string),
		float32(temperature),
		int(maxTokens),
		float32(topP),
		float32(frequencyPenalty),
		input.UserMessage,
	)
	fmt.Println("AI 响应:", response)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"chat_id": input.ChatID, "AI_response": response}})

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
