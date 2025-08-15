package config

import (
	"os"
)

type AIModel struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Logo     string `json:"logo"`
	Ability  string `json:"ability"`
	Price    string `json:"price"`
	API_KEY  string `json:"api_key"`
	BASE_URL string `json:"base_url"`
}

var AIModels = []AIModel{
	{
		Key:      "claude-sonnet-4-20250514",
		Name:     "Claude-Sonnet-4",
		Logo:     "Claude.svg",
		Ability:  "image",
		Price:    "100",
		API_KEY:  "AI_API_KEY",
		BASE_URL: "AI_BASE_URL",
	},
	{
		Key:      "gpt-5",
		Name:     "ChatGPT-5",
		Logo:     "ChatGPT.svg",
		Ability:  "text",
		Price:    "100",
		API_KEY:  "AI_API_KEY",
		BASE_URL: "AI_BASE_URL",
	},
	{
		Key:      "gemini-2.5-flash-nothink",
		Name:     "Gemini-2.5-Flash",
		Logo:     "Gemini.svg",
		Ability:  "",
		Price:    "100",
		API_KEY:  "AI_API_KEY",
		BASE_URL: "AI_BASE_URL",
	},
	{
		Key:      "gemini-2.5-pro-preview-05-06",
		Name:     "Gemini-2.5-Pro",
		Logo:     "Gemini.svg",
		Ability:  "image,audio,video",
		Price:    "100",
		API_KEY:  "AI_API_KEY",
		BASE_URL: "AI_BASE_URL",
	},
	{
		Key:      "qwen-vl-plus-2025-05-07",
		Name:     "Qwen-VL-Plus",
		Logo:     "Qwen.svg",
		Ability:  "image,audio,video",
		Price:    "100",
		API_KEY:  "DASHSCOPE_API_KEY",
		BASE_URL: "DASHSCOPE_BASE_URL",
	},
}

// 根据name获取key
func GetModelKey(name string) string {
	for _, model := range AIModels {
		if model.Name == name {
			return model.Key
		}
	}
	return ""
}

// 根据name获取API_KEY和BASE_URL
func GetAIConfig(name string) (string, string) {
	for _, model := range AIModels {
		if model.Name == name {
			apiKey := os.Getenv(model.API_KEY)
			baseURL := os.Getenv(model.BASE_URL)
			return apiKey, baseURL
		}
	}
	return "", ""
}
