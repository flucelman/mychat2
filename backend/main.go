package main

import (
	"backend/config"
	"backend/router"
)

func main() {
	// 初始化配置
	config.InitConfig()
	// 初始化路由
	r := router.SetupRouter()
	// 启动服务
	r.Run("0.0.0.0:" + config.AppConfig.App.Port)
}
