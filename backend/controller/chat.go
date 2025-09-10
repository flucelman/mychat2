package controller

import (
	"backend/config"
	"backend/global"
	"backend/models"
	"backend/utils"
	"backend/utils/grpc"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
查询聊天记录表
通过user_id查询聊天记录
*/
func GetChatHistory(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	chatHistory := []models.ChatHistory{}

	// 在查询时直接按 updated_at 降序排列（最新的在前）
	if err := global.DB.Where("user_id = ?", userID).Order("updated_at desc").Find(&chatHistory).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type ChatHistoryResponse struct {
		ChatID    string    `json:"chat_id"`
		Title     string    `json:"title"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	chatHistoryResponse := []ChatHistoryResponse{}
	for _, chat := range chatHistory {
		chatHistoryResponse = append(chatHistoryResponse, ChatHistoryResponse{
			ChatID:    chat.ChatID,
			Title:     chat.Title,
			UpdatedAt: chat.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"data": chatHistoryResponse})
}

/*
查询消息表
通过chat_id查询消息

	返回：{
		message_id: string,
		role: string,
		content: string,
		model: string,
		created_at: time.Time,
	}
*/
func GetChatMessage(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	chatID := ctx.Param("chat_id")

	// 从查询参数获取分页信息，设置默认值
	page := 1
	pageSize := 20

	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}
	message := []models.Message{}
	if err := global.DB.Where("user_id = ?", userID).Where("chat_id = ?", chatID).
		Order("created_at ASC"). // 按创建时间升序排列（最早的在前）
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&message).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 查询文件
	file := []models.File{}
	if err := global.DB.Where("user_id = ?", userID).Where("chat_id = ?", chatID).Find(&file).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 定义统一的响应结构
	type CombinedResponse struct {
		MessageID   string    `json:"message_id,omitempty"`
		FileID      string    `json:"file_id,omitempty"`
		Role        string    `json:"role"`
		Content     string    `json:"content,omitempty"`
		Model       string    `json:"model,omitempty"`
		FileURL     string    `json:"file_url,omitempty"`
		FileSize    int64     `json:"file_size,omitempty"`
		FileName    string    `json:"file_name,omitempty"`
		FileType    string    `json:"file_type,omitempty"`
		FileContent string    `json:"file_content,omitempty"`
		CreatedAt   time.Time `json:"created_at"`
	}

	var combinedResponse []CombinedResponse

	// 添加消息到结果中
	for _, msg := range message {
		combinedResponse = append(combinedResponse, CombinedResponse{
			MessageID: msg.MessageID,
			Role:      msg.Role,
			Content:   msg.Content,
			Model:     msg.Model,
			CreatedAt: msg.CreatedAt,
		})
	}

	// 添加文件到结果中
	for _, f := range file {
		combinedResponse = append(combinedResponse, CombinedResponse{
			FileID:      f.FileID,
			Role:        "file",
			FileURL:     f.FileURL,
			FileSize:    f.FileSize,
			FileName:    f.FileName,
			FileType:    f.FileType,
			FileContent: f.FileContent,
			CreatedAt:   f.CreatedAt,
		})
	}

	// 按照 CreatedAt 排序（升序）
	for i := 0; i < len(combinedResponse)-1; i++ {
		for j := 0; j < len(combinedResponse)-1-i; j++ {
			if combinedResponse[j].CreatedAt.After(combinedResponse[j+1].CreatedAt) {
				combinedResponse[j], combinedResponse[j+1] = combinedResponse[j+1], combinedResponse[j]
			}
		}
	}

	ctx.JSON(http.StatusOK, combinedResponse)
}

// 查询模型列表
func GetModelList(ctx *gin.Context) {
	type ModelList struct {
		Key     string `json:"key"`
		Name    string `json:"name"`
		Logo    string `json:"logo"`
		Ability string `json:"ability"`
		Price   string `json:"price"`
	}
	modelList := []ModelList{}
	for _, model := range config.AIModels {
		modelList = append(modelList, ModelList{
			Key:     model.Key,
			Name:    model.Name,
			Logo:    model.Logo,
			Ability: model.Ability,
			Price:   model.Price,
		})
	}
	ctx.JSON(http.StatusOK, modelList)
}

/*
添加对话信息并通过SSE流式返回AI响应
如果是空聊天记录，则创建新的聊天记录
如果非空，则添加对话信息
*/
func AddChatMessage(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var input struct {
		ChatID         string           `json:"chat_id"`
		AIConfig       map[string]any   `json:"AI_config"`
		MessageHistory []map[string]any `json:"message_history"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 计算MessageHistory里面有多少个user消息
	userCount := 0
	for _, msg := range input.MessageHistory {
		if msg["role"] == "user" {
			userCount++
		}
	}
	// 如果只有一个user消息，则检查是否需要创建新的聊天记录
	if userCount == 1 {
		// 先检查聊天记录是否已存在
		var existingChat models.ChatHistory
		err := global.DB.Where("chat_id = ? AND user_id = ?", input.ChatID, userID).First(&existingChat).Error

		if err != nil && err == gorm.ErrRecordNotFound {
			// 聊天记录不存在，创建新的聊天记录
			fmt.Println("创建新的聊天记录", input.MessageHistory)
			// 保存system prompt信息到数据库
			systemPrompt := input.MessageHistory[0]["content"]
			systemPromptID := uuid.New().String()
			utils.SaveDB(systemPromptID, userID, input.ChatID, "system", systemPrompt.(string), "system")
			input.MessageHistory[0]["message_id"] = systemPromptID
			chat := models.ChatHistory{
				ChatID: input.ChatID,
				UserID: userID,
				Title: func() string {
					// 安全检查：确保MessageHistory有足够的元素
					if len(input.MessageHistory) < 2 {
						return "新对话"
					}

					// 找到第一条用户消息作为标题
					for _, msg := range input.MessageHistory {
						if role, ok := msg["role"].(string); ok && role == "user" {
							if content, ok := msg["content"].(string); ok && content != "" {
								if utf8.RuneCountInString(content) <= 30 {
									return content
								}
								// 安全截取UTF-8字符串，确保不会截断中文字符
								runes := []rune(content)
								if len(runes) > 32 {
									return string(runes[:32]) + "..."
								}
								return content
							}
						}
					}

					// 如果都失败了，返回默认标题
					fmt.Println("都失败了，返回默认标题")
					return "新对话"
				}(),
			}
			if err := global.DB.Create(&chat).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else if err != nil {
			// 其他数据库错误
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			// 聊天记录已存在，不需要创建
			fmt.Println("聊天记录已存在，重发消息")
		}
	}
	// 2. 添加用户信息到数据库
	if len(input.MessageHistory) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "消息历史为空"})
		fmt.Println("消息历史为空")
		return
	}

	// 找到最后一条用户消息（role为user的消息）
	var lastUserMessage map[string]any
	for i := len(input.MessageHistory) - 1; i >= 0; i-- {
		if role, ok := input.MessageHistory[i]["role"].(string); ok && role == "user" {
			lastUserMessage = input.MessageHistory[i]
			break
		}
	}

	if lastUserMessage == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "未找到用户消息"})
		return
	}

	messageID, ok := lastUserMessage["message_id"].(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息ID类型"})
		return
	}
	content, ok := lastUserMessage["content"].(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息内容类型"})
		return
	}
	model, ok := input.AIConfig["model"].(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的模型类型"})
		return
	}
	// 保存用户信息到数据库
	go utils.SaveDB(messageID, userID, input.ChatID, "user", content, model)

	// 检查是否开启的联网搜索
	onlineSearch := input.AIConfig["online_search"]
	if onlineSearch == true {
		onlineSearchResponse, err := grpc.OnlineSearch(content)
		if err != nil {
			fmt.Println("联网搜索失败", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		searchMessageID := uuid.New().String()
		go utils.SaveDB(searchMessageID, userID, input.ChatID, "search", onlineSearchResponse, model)
		fmt.Println("联网搜索响应：", onlineSearchResponse)
		ctx.SSEvent("search", gin.H{"search_result": onlineSearchResponse, "message_id": searchMessageID})
		input.MessageHistory = append(input.MessageHistory, map[string]any{
			"role":       "search",
			"content":    onlineSearchResponse,
			"message_id": searchMessageID,
		})
		ctx.Writer.Flush()
	}

	// 3. 设置SSE响应头
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Cache-Control")

	// 5. 发送初始响应（包含chat_id）
	assistantMessageID := uuid.New().String()
	ctx.SSEvent("start", gin.H{"chat_id": input.ChatID, "assistant_message_id": assistantMessageID})
	ctx.Writer.Flush()

	// 7. 创建通道和上下文
	answerCh := make(chan string)
	answerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 8. 启动AI流式响应
	modelKey := config.GetModelKey(model) // 使用之前已经验证过的model变量
	apiKey, baseURL := config.GetAIConfig(model)
	if input.AIConfig["planExecutor"] == true {
		fmt.Println("启动计划执行器")
		go utils.PlanExecutorResponse(answerCtx, answerCh, apiKey, baseURL, modelKey, input.AIConfig, input.MessageHistory)
	} else {
		fmt.Println("启动AI流式响应")
		go utils.AIStreamResponse(answerCtx, answerCh, apiKey, baseURL, modelKey, input.AIConfig, input.MessageHistory)
	}

	// 9. 处理流式响应并通过SSE发送
	fullResponse := ""
	for {
		select {
		case content, ok := <-answerCh:
			if !ok {
				// 通道关闭，AI响应结束
				// 保存完整的AI响应到数据库
				fmt.Println("助手消息ID：", assistantMessageID)
				saveResponse := utils.SaveDB(assistantMessageID, userID, input.ChatID, "assistant", fullResponse, model)

				if saveResponse != "success" {
					ctx.SSEvent("error", gin.H{"error": saveResponse})
				} else {
					ctx.SSEvent("end", gin.H{"message": "AI响应完成", "message_id": assistantMessageID})
				}
				return
			}

			// 检查是否是错误消息
			if len(content) > 6 && content[:6] == "ERROR:" {
				// 发送错误信息并断开连接
				fmt.Println("gRPC服务错误：", content)
				ctx.SSEvent("error", gin.H{"error": content[6:]}) // 去掉ERROR:前缀
				ctx.Writer.Flush()
				cancel()
				return
			}

			// 发送内容片段
			ctx.SSEvent("content", gin.H{"content": content})
			fullResponse += content
			ctx.Writer.Flush()

		case <-ctx.Request.Context().Done():
			// 客户端断开连接
			fmt.Println("客户端断开连接，停止AI响应")
			if fullResponse == "" {
				fmt.Println("助手消息为空，不保存到数据库")
				cancel()
				return
			}
			fmt.Println("助手消息ID：", assistantMessageID)
			saveResponse := utils.SaveDB(assistantMessageID, userID, input.ChatID, "assistant", fullResponse, model)
			if saveResponse != "success" {
				ctx.SSEvent("error", gin.H{"error": saveResponse})
			} else {
				ctx.SSEvent("end", gin.H{"message": "AI响应完成", "message_id": assistantMessageID})
			}
			cancel()
			return
		}
	}
}

// 删除所有历史记录
func DeleteAllHistory(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if err := global.DB.Where("user_id = ?", userID).Delete(&models.ChatHistory{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := global.DB.Where("user_id = ?", userID).Delete(&models.Message{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "所有历史记录已删除"})
}

// 删除单个聊天记录
func DeleteSingleHistory(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var input struct {
		ChatID string `json:"chat_id"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := global.DB.Where("user_id = ?", userID).Where("chat_id = ?", input.ChatID).Delete(&models.ChatHistory{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := global.DB.Where("user_id = ?", userID).Where("chat_id = ?", input.ChatID).Delete(&models.Message{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "聊天记录已删除"})
}

// 删除单个消息
func DeleteSingleMessage(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var input struct {
		MessageID string `json:"message_id"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := global.DB.Where("user_id = ?", userID).Where("message_id = ?", input.MessageID).Delete(&models.Message{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "消息已删除"})
}

// 重发消息
func ResendMessage(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	var input struct {
		ChatID    string `json:"chat_id"`
		MessageID string `json:"message_id"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. 先查找要重发的消息
	var targetMessage models.Message
	if err := global.DB.Where("user_id = ? AND chat_id = ? AND message_id = ?", userID, input.ChatID, input.MessageID).First(&targetMessage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "消息不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 2. 获取该聊天下的所有消息，按创建时间排序
	var allMessages []models.Message
	if err := global.DB.Where("user_id = ? AND chat_id = ?", userID, input.ChatID).
		Order("created_at ASC").Find(&allMessages).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. 找到目标消息的索引位置
	targetIndex := -1
	for i, msg := range allMessages {
		if msg.MessageID == input.MessageID {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "消息不存在"})
		return
	}

	// 4. 确定要删除的消息范围
	// 删除目标消息及其之后的所有消息，还要删除目标消息的上一个消息
	deleteStartIndex := targetIndex
	if targetIndex > 0 {
		deleteStartIndex = targetIndex - 1 // 包括上一个消息
	}

	// 5. 收集要删除的消息ID
	var messageIDsToDelete []string
	for i := deleteStartIndex; i < len(allMessages); i++ {
		messageIDsToDelete = append(messageIDsToDelete, allMessages[i].MessageID)
	}

	// 6. 批量删除消息
	if len(messageIDsToDelete) > 0 {
		if err := global.DB.Where("user_id = ? AND chat_id = ? AND message_id IN ?",
			userID, input.ChatID, messageIDsToDelete).Delete(&models.Message{}).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除消息失败: " + err.Error()})
			return
		}
		fmt.Printf("删除了 %d 条消息: %v\n", len(messageIDsToDelete), messageIDsToDelete)
	}

	// 7. 返回成功响应，包含删除的消息数量
	ctx.JSON(http.StatusOK, gin.H{
		"success":             true,
		"message":             "消息已重发",
		"deleted_count":       len(messageIDsToDelete),
		"deleted_message_ids": messageIDsToDelete,
	})
}
