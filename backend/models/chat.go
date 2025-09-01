package models

import "time"

// 聊天记录表
type ChatHistory struct {
	ChatID    string    `json:"chat_id" gorm:"primaryKey;type:varchar(64)"`
	UserID    string    `json:"user_id" gorm:"type:varchar(64);index;not null"`
	Title     string    `json:"title" gorm:"type:varchar(64)"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// 消息表
type Message struct {
	MessageID string    `json:"message_id" gorm:"primaryKey;type:varchar(64)"`
	ChatID    string    `json:"chat_id" gorm:"type:varchar(64);index;not null"`
	UserID    string    `json:"user_id" gorm:"type:varchar(64);index;not null"`
	FileID    string    `json:"file_id" gorm:"type:varchar(64);index;not null"`
	Role      string    `json:"role" gorm:"type:varchar(32);not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Model     string    `json:"model" gorm:"type:varchar(32);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamp"`
}

type File struct {
	FileID      string    `json:"file_id" gorm:"primaryKey;type:varchar(64)"`
	UserID      string    `json:"user_id" gorm:"type:varchar(64);index;not null"`
	ChatID      string    `json:"chat_id" gorm:"type:varchar(64);index;not null"`
	FileName    string    `json:"file_name" gorm:"type:varchar(255);not null"`
	FileSize    int64     `json:"file_size" gorm:"not null"`
	FileURL     string    `json:"file_url" gorm:"type:varchar(255);not null"`
	FileType    string    `json:"file_type" gorm:"type:varchar(32);default:file"`
	FileContent string    `json:"file_content" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamp"`
}
