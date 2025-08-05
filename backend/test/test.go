package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

// 调用AI接口，使用channel进行流式输出
func AIResponse(ctx context.Context, model string, temperature float32, max_tokens int, top_p float32, frequency_penalty float32, content string, ch chan<- string) {
	defer close(ch) // 确保在函数结束时关闭channel

	godotenv.Load()
	config := openai.DefaultConfig(os.Getenv("AI_API_KEY"))

	// 如果设置了自定义 AI_URL，则使用自定义端点
	if aiURL := os.Getenv("AI_BASE_URL"); aiURL != "" {
		config.BaseURL = aiURL
	}

	client := openai.NewClientWithConfig(config)

	stream, err := client.CreateChatCompletionStream(
		ctx, // 使用传入的context
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
			Stream:           true,
		},
	)
	if err != nil {
		ch <- fmt.Sprintf("错误: %s", err.Error())
		return
	}
	defer stream.Close()

	for {
		select {
		case <-ctx.Done():
			// 当context被取消时，发送取消信息并返回
			ch <- fmt.Sprintf("\n[流输出已被取消: %s]", ctx.Err().Error())
			return
		default:
			response, err := stream.Recv()
			if err != nil {
				// 检查是否是正常的流结束信号
				if err == io.EOF {
					fmt.Println("\n[AI流式响应正常结束]")
					return
				}
				// 其他错误才发送到channel
				ch <- fmt.Sprintf("\n[流输出错误: %s]", err.Error())
				return
			}

			if response.Usage != nil {
				fmt.Printf("Usage: PromptTokens=%d, CompletionTokens=%d, TotalTokens=%d\n",
					response.Usage.PromptTokens,
					response.Usage.CompletionTokens,
					response.Usage.TotalTokens)
			}

			if len(response.Choices) > 0 {
				ch <- response.Choices[0].Delta.Content
			}
		}
	}
}

func main() {
	input := "请写一篇关于人工智能发展历史的长文章，包含详细的时间线和重要事件"
	ch := make(chan string, 50)

	// 创建一个3秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 启动goroutine执行AI响应
	go AIResponse(ctx, "gpt-4o-mini", 0.5, 1000, 0.5, 0.5, input, ch)

	fmt.Print("AI回复: ")
	// 实时接收并打印响应内容
	for content := range ch {
		fmt.Print(content)
	}
	fmt.Println() // 换行
	fmt.Println("程序结束")
}
