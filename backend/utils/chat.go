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

	// 只更新已存在的聊天记录的 UpdatedAt 字段
	// 使用 Updates 方法直接更新，避免先查询再保存的问题
	result := global.DB.Model(&models.ChatHistory{}).
		Where("chat_id = ?", chatID).
		Update("updated_at", time.Now())

	// 如果没有找到记录，不需要报错，因为可能是第一次创建聊天时的system消息
	// ChatHistory 会在 controller/chat.go 的 AddChatMessage 函数中创建
	if result.Error != nil {
		// 记录错误但不返回错误，因为消息已经保存成功
		// 可以考虑记录日志
		return "success"
	}

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

func PlanExecutorResponse(
	ctx context.Context,
	answerCh chan<- string,
	apiKey,
	baseURL,
	model string,
	aiConfig map[string]any,
	message_history []map[string]any) {
	grpc.PlanExecutorResponse(ctx, answerCh, model, apiKey, baseURL, aiConfig, message_history)
}
