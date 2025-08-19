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
	aiConfig map[string]any,
	message_history []map[string]any) {

	grpc.AIStreamResponse(ctx, answerCh, model, apiKey, baseURL, aiConfig, message_history)
}
