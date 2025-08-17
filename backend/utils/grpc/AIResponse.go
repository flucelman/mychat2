package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
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
	temperature float32,
	maxTokens int,
	topP float32, frequencyPenalty float32, messageHistory []map[string]any) string {
	addr := os.Getenv("GRPC_AIRESPONSE")
	fmt.Println("addr:", addr)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	client := AIResponse.NewGreeterClient(conn)

	// 转换消息历史格式
	messages := make([]*AIResponse.Message, 0, len(messageHistory))
	for _, msg := range messageHistory {
		message := &AIResponse.Message{
			Role:    msg["role"].(string),
			Content: msg["content"].(string),
		}
		// 可选字段
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

	fmt.Println("messages格式", messages)

	// 创建请求
	request := &AIResponse.AIStreamRequest{
		Model:            model,
		ApiKey:           apiKey,
		BaseUrl:          baseUrl,
		Temperature:      temperature,
		MaxTokens:        int32(maxTokens),
		TopP:             topP,
		FrequencyPenalty: frequencyPenalty,
		MessageHistory:   messages,
	}

	// 调用服务
	stream, err := client.AIResponse(ctx, request)
	if err != nil {
		log.Fatalf("调用服务失败: %v", err)
	}

	fmt.Println("AI 响应:")
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// 流结束
			break
		}
		if err != nil {
			log.Fatalf("接收响应失败: %v", err)
		}
		// 打印每个流式响应片段
		fmt.Print(response.GetContent())
		answerCh <- response.GetContent()
	}
	fmt.Println() // 最后打印一个换行
	close(answerCh)
	return ""
}
