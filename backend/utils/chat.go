package utils

import (
	"backend/global"
	"backend/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

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

func AIResponse(Model, Content string) string {
	url := os.Getenv("AI_URL") + "/chat/completions"
	jsonStr := `{
		"model": "` + Model + `",
		"messages": [
		  {
			"role": "user",
			"content": "` + Content + `"
		  }
		],
		"temperature": 0.8,
		"max_tokens": 1024,
		"top_p": 1,
		"frequency_penalty": 0,
		"presence_penalty": 0,
		"stream": false
	  }`
	payload := strings.NewReader(strings.ReplaceAll(jsonStr, "\n", ""))
	request, err := http.NewRequest("POST", url, payload)
	fmt.Println("请求体request:", request)
	if err != nil {
		return err.Error()
	}
	request.Header.Add("Authorization", "Bearer "+os.Getenv("AI_API_KEY"))
	request.Header.Add("Content-Type", "application/json")
	fmt.Println("开始发送请求")
	client := &http.Client{}
	response, err := client.Do(request)
	fmt.Println("请求响应response:", response)
	if err != nil {
		return err.Error()
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	return string(body)
}
