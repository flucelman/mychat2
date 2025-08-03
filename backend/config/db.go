package config

import (
	"log"
	"os"
	"time"

	"backend/global"
	"backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	db_password := os.Getenv("DB_PASSWORD")
	dsn := AppConfig.Database.Username + ":" + db_password + "@tcp(" + AppConfig.Database.Host + ":" + AppConfig.Database.Port + ")/" + AppConfig.Database.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 创建表
	var tables = []any{
		&models.User{},
		&models.UserToken{},
		&models.ChatHistory{},
		&models.Message{},
	}

	// 自动迁移数据库
	if err := db.AutoMigrate(tables...); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	global.DB = db
}
