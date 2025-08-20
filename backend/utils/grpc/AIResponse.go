package grpc

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"backend/utils/grpc/pb/AIResponse"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AIStreamResponse(
	ctx context.Context,
	answerCh chan<- string,
	model string,
	apiKey string,
	baseUrl string,
	aiConfig map[string]any,
	messageHistory []map[string]any) string {
	addr := os.Getenv("GRPC_AIRESPONSE")
	fmt.Println("addr:", addr)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接失败: ", err)
		return ""
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	client := AIResponse.NewGreeterClient(conn)

	// 转换消息历史格式
	messages := make([]*AIResponse.Message, 0, len(messageHistory))
	for _, msg := range messageHistory {
		message := &AIResponse.Message{
			Role: msg["role"].(string),
		}
		// 可选字段
		if content, ok := msg["content"].(string); ok {
			message.Content = content
		}
		if fileType, ok := msg["file_type"].(string); ok {
			message.FileType = fileType
		}
		if fileUrl, ok := msg["file_url"].(string); ok {
			message.FileUrl = fileUrl
		}
		if fileContent, ok := msg["file_content"].(string); ok {
			message.FileContent = fileContent
		}
		messages = append(messages, message)
	}

	// 创建请求
	request := &AIResponse.AIStreamRequest{
		Model:            model,
		ApiKey:           apiKey,
		BaseUrl:          baseUrl,
		Temperature:      float32(aiConfig["temperature"].(float64)),
		MaxTokens:        int32(aiConfig["max_tokens"].(float64)),
		TopP:             float32(aiConfig["top_p"].(float64)),
		FrequencyPenalty: float32(aiConfig["frequency_penalty"].(float64)),
		MessageHistory:   messages,
	}

	// 调用服务
	stream, err := client.AIResponse(ctx, request)
	if err != nil {
		fmt.Println("调用服务失败: ", err)
		return ""
	}
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// 流结束
			break
		}
		if err != nil {
			fmt.Println("接收响应失败: ", err)
			return ""
		}
		// 打印每个流式响应片段
		answerCh <- response.GetContent()
	}
	close(answerCh)
	return ""
}
