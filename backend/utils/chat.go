package utils

import (
	"backend/global"
	"backend/models"
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

// 保存对话信息到数据库
func SaveDB(userID, chatID, Role, Content, Model string) string {
	message := models.Message{
		MessageID: uuid.New().String(),
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
	return "success"
}

// 调用AI接口
func AIResponse(model string, temperature float32, max_tokens int, top_p float32, frequency_penalty float32, content string) string {
	config := openai.DefaultConfig(os.Getenv("AI_API_KEY"))

	// 如果设置了自定义 AI_URL，则使用自定义端点
	if aiURL := os.Getenv("AI_BASE_URL"); aiURL != "" {
		config.BaseURL = aiURL
	}

	client := openai.NewClientWithConfig(config)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
			Temperature:      temperature,
			MaxTokens:        max_tokens,
			TopP:             top_p,
			FrequencyPenalty: frequency_penalty,
		},
	)

	if err != nil {
		return err.Error()
	}

	return resp.Choices[0].Message.Content
}
