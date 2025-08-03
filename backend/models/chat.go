package models

import "time"

// 聊天记录表
type ChatHistory struct {
	ChatID    string    `json:"chat_id" gorm:"primaryKey;type:varchar(64)"`
	UserID    string    `json:"user_id" gorm:"type:varchar(64);index;not null"`
	Title     string    `json:"title" gorm:"type:varchar(32)"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// 消息表
type Message struct {
	MessageID string    `json:"message_id" gorm:"primaryKey;type:varchar(64)"`
	ChatID    string    `json:"chat_id" gorm:"type:varchar(64);index;not null"`
	UserID    string    `json:"user_id" gorm:"type:varchar(64);index;not null"`
	Role      string    `json:"role" gorm:"type:varchar(32);not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Model     string    `json:"model" gorm:"type:varchar(32);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamp"`
}
