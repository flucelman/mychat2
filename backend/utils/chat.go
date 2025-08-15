package utils

import (
	"backend/global"
	"backend/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
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

// 调用AI接口
func AIStreamResponse(ctx context.Context, answerCh chan<- string, apiKey, baseURL, model string, temperature float32, max_tokens int, top_p float32, frequency_penalty float32, message_history []map[string]any) {

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL

	// 构建消息数组
	messages := []openai.ChatCompletionMessage{}

	// 添加历史消息
	for _, msg := range message_history {
		// 检查 role 字段
		roleInterface, ok := msg["role"]
		if !ok || roleInterface == nil {
			continue // 跳过没有 role 字段的消息
		}
		role, ok := roleInterface.(string)
		if !ok {
			continue // 跳过 role 不是字符串的消息
		}

		// 检查 content 字段
		contentInterface, ok := msg["content"]
		// 先处理文件类消息，避免对 nil content 做类型断言
		if role == "file" {
			fileMsgs := handleFile(msg)
			if len(fileMsgs) > 0 {
				messages = append(messages, fileMsgs...)
			} else {
				fmt.Println("file 消息未生成有效内容，跳过")
			}
			continue
		}

		// 处理非文件类消息，需要安全地读取 content
		if !ok || contentInterface == nil {
			fmt.Println("content 缺失或为 nil，跳过")
			continue
		}
		content, ok := contentInterface.(string)
		if !ok || strings.TrimSpace(content) == "" {
			fmt.Println("content 不是字符串或为空，跳过")
			continue
		}
		var openaiRole string
		switch role {
		case "user":
			openaiRole = openai.ChatMessageRoleUser
		case "assistant":
			openaiRole = openai.ChatMessageRoleAssistant
		case "system":
			openaiRole = openai.ChatMessageRoleSystem
		default:
			openaiRole = openai.ChatMessageRoleUser // 默认为用户角色
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openaiRole,
			Content: content,
		})
	}
	jsonData, _ := json.Marshal(messages)
	fmt.Println("messages设置：", string(jsonData))
	client := openai.NewClientWithConfig(config)
	stream, err := client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:            model,
			Messages:         messages,
			Temperature:      temperature,
			MaxTokens:        max_tokens,
			TopP:             top_p,
			FrequencyPenalty: frequency_penalty,
			Stream:           true,
		},
	)
	if err != nil {
		fmt.Printf("创建流失败: %v\n", err)
		close(answerCh)
		return
	}
	defer stream.Close() // 确保在函数结束时关闭stream

	// 修复：在循环中持续检查取消信号
	for {
		response, err := stream.Recv()
		select {
		case <-ctx.Done(): // 用户手动取消或超时
			fmt.Println("收到取消信号，停止AI响应")
			if response.Usage != nil {
				fmt.Printf("Usage: PromptTokens=%d, CompletionTokens=%d, TotalTokens=%d\n",
					response.Usage.PromptTokens,
					response.Usage.CompletionTokens,
					response.Usage.TotalTokens)
			} else {
				fmt.Println("Usage: nil", response)
			}
			stream.Close()
			close(answerCh)
			return
		default:
			if err != nil {
				// 检查是否是正常的流结束信号
				if err == io.EOF {
					fmt.Println("AI流式响应正常结束")
					close(answerCh)
					return
				}
				// 其他错误才打印错误信息
				fmt.Printf("接收流数据失败: %v\n", err)
				close(answerCh)
				return
			}
			if response.Usage != nil {
				fmt.Printf("Usage: PromptTokens=%d, CompletionTokens=%d, TotalTokens=%d\n",
					response.Usage.PromptTokens,
					response.Usage.CompletionTokens,
					response.Usage.TotalTokens)
			}

			if len(response.Choices) > 0 {
				answerCh <- response.Choices[0].Delta.Content
			}
		}
	}
}

func handleFile(msg map[string]any) []openai.ChatCompletionMessage {
	var messages []openai.ChatCompletionMessage
	switch msg["file_type"] {
	case "image":
		fmt.Println("image", msg["file_url"])
		messages = append(messages, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser,
			MultiContent: []openai.ChatMessagePart{
				{
					Type: openai.ChatMessagePartTypeImageURL,
					ImageURL: &openai.ChatMessageImageURL{
						URL: msg["file_url"].(string),
					},
				},
			},
		})
		return messages
	case "file":
		fmt.Println("开始解析file")
		fileContent := "parse file failed"
		// 认为空字符串等同于没有内容，需要触发解析
		fmt.Println("msg的内容:", msg)
		if s, ok := msg["file_content"].(string); ok && strings.TrimSpace(s) != "" {
			fmt.Println("file_content不为空")
			fileContent = s
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: "file_content: " + fileContent,
		})
	case "video":

		messages = append(messages, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser,
			MultiContent: []openai.ChatMessagePart{
				{
					Type: openai.ChatMessagePartTypeVideoURL,
					VideoURL: map[string]any{
						"url": msg["file_url"].(string),
					},
				},
			},
		})
		return messages
	}
	return messages
}
