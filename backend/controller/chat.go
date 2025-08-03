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
func AddChatMessage(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var input struct {
		ChatID  string `json:"chat_id"`
		Role    string `json:"role"`
		Content string `json:"content"`
		Model   string `json:"model"`
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
				if len(input.Content) <= 10 {
					return input.Content
				}
				return input.Content[:10]
			}(),
		}
		if err := global.DB.Create(&chat).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	// 2. 添加对话信息
	response := utils.SaveDB(userID, input.ChatID, input.Role, input.Content, input.Model)
	if response != "success" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": response})
		return
	}

	// 3. 判断是用户还是助手
	fmt.Printf("收到的 Role: %s\n", input.Role)
	if input.Role == "user" {
		fmt.Println("开始调用 AI 接口...")
		response = utils.AIResponse(input.Model, input.Content)
		fmt.Println("AI 响应:", response)
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "chat_id": input.ChatID})
}
