package models

import (
	"time"
)

type User struct {
	Id         int        `json:"id"`
	TelegramId int64      `json:"telegram_id"`
	Name       string     `json:"name"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	CreatedAt  *time.Time `json:"created_at"`
}

type Token struct {
	Id        int        `json:"id"`
	UserId    int        `json:"user_id"`
	Platform  string     `json:"platform"`
	DeviceId  string     `json:"device_id"`
	Value     string     `json:"value"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type Note struct {
	Id        int64      `json:"id"`
	UserId    int        `json:"user_id"`
	Text      string     `json:"text"`
	FileId    string     `json:"file_id"`
	FileType  byte       `json:"file_type"`
	CreatedAt *time.Time `json:"created_at"`
}
