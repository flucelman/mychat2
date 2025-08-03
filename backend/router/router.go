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
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
	}

	api := r.Group("/api")
	// 中间件，验证token
	api.Use(middlewares.AuthMiddleware())
	{
		chat := api.Group("/chat")
		{
			chat.POST("/history", controller.GetChatHistory)
			chat.POST("/message/:chat_id", controller.GetChatMessage)
			chat.POST("/add_message", controller.AddChatMessage)
		}
	}

	return r
}
