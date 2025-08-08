package config

type AIModel struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	Ability string `json:"ability"`
	Price   string `json:"price"`
}

var AIModels = []AIModel{
	{
		Key:     "claude-sonnet-4-20250514",
		Name:    "Claude-Sonnet-4",
		Ability: "image",
		Price:   "100",
	},
	{
		Key:     "gpt-5",
		Name:    "ChatGPT-5",
		Ability: "",
		Price:   "100",
	},
	{
		Key:     "gemini-2.5-flash-nothink",
		Name:    "Gemini-2.5-Flash",
		Ability: "",
		Price:   "100",
	},
	{
		Key:     "gemini-2.5-pro-preview-05-06",
		Name:    "Gemini-2.5-Pro",
		Ability: "image,audio,video",
		Price:   "100",
	},
}
