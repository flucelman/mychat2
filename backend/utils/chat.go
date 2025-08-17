package utils

import (
	"backend/global"
	"backend/models"
	"backend/utils/grpc"
	"context"
	"time"
)

// 保存对话信息到数据库
func SaveDB(MessageID, userID, chatID, Role, Content, Model string) string {
	message := models.Message{
		MessageID: MessageID,
		ChatID:    chatID,
		UserID:    userID,
		Role:      Role,
		Content:   Content,
		Model:     Model,
	}
	// 如果role为user，则添加对话信息
	if Role == "user" {
		message.Role = "user"
	} else {
		// 如果role为assistant，则添加对话信息
		message.Role = "assistant"
	}
	message.Content = Content
	if err := global.DB.Create(&message).Error; err != nil {
		return err.Error()
	}
	chatHistory := models.ChatHistory{}
	global.DB.Model(&models.ChatHistory{}).Where("chat_id = ?", chatID).First(&chatHistory)
	chatHistory.UpdatedAt = time.Now()
	global.DB.Save(&chatHistory)
	return "success"
}

// 调用AI接口
func AIStreamResponse(
	ctx context.Context,
	answerCh chan<- string,
	apiKey,
	baseURL,
	model string,
	temperature float32,
	max_tokens int,
	top_p float32,
	frequency_penalty float32,
	message_history []map[string]any) {

	answerCh <- grpc.AIStreamResponse(ctx, answerCh, model, apiKey, baseURL, temperature, max_tokens, top_p, frequency_penalty, message_history)
}
