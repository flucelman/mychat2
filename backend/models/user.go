package models

import (
	"time"
)

type User struct {
	UserID    string    `json:"user_id" gorm:"primaryKey;unique;type:varchar(64)"`
	Email     string    `json:"email" gorm:"unique"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type UserToken struct {
	UserID string `json:"user_id" gorm:"primaryKey;unique;type:varchar(64)"`
	Token  string `json:"token"`
}
