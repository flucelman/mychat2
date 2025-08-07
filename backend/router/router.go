package router

import (
	"backend/controller"
	"backend/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 配置跨域
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
		auth.POST("/checkToken", controller.CheckToken)
	}

	api := r.Group("/api")
	// 中间件，验证token
	api.Use(middlewares.AuthMiddleware())
	{
		chat := api.Group("/chat")
		{
			chat.GET("/getChatHistory", controller.GetChatHistory)
			chat.GET("/getChatMessage/:chat_id", controller.GetChatMessage)
			chat.POST("/addChatMessage", controller.AddChatMessage)
			chat.DELETE("/deleteAllHistory", controller.DeleteAllHistory)
			chat.DELETE("/deleteSingleHistory/:chat_id", controller.DeleteSingleHistory)
		}
	}

	return r
}
