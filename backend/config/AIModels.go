package config

type AIModel struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	Logo    string `json:"logo"`
	Ability string `json:"ability"`
	Price   string `json:"price"`
}

var AIModels = []AIModel{
	{
		Key:     "claude-sonnet-4-20250514",
		Name:    "Claude-Sonnet-4",
		Logo:    "Claude.svg",
		Ability: "image",
		Price:   "100",
	},
	{
		Key:     "gpt-5",
		Name:    "ChatGPT-5",
		Logo:    "ChatGPT.svg",
		Ability: "text",
		Price:   "100",
	},
	{
		Key:     "gemini-2.5-flash-nothink",
		Name:    "Gemini-2.5-Flash",
		Logo:    "Gemini.svg",
		Ability: "",
		Price:   "100",
	},
	{
		Key:     "gemini-2.5-pro-preview-05-06",
		Name:    "Gemini-2.5-Pro",
		Logo:    "Gemini.svg",
		Ability: "image,audio,video",
		Price:   "100",
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
