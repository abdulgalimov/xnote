package models

import (
	"time"
)

// User структура юзера
type User struct {
	ID         int        `json:"id"`
	TelegramID int64      `json:"telegram_id"`
	Name       string     `json:"name"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	CreatedAt  *time.Time `json:"created_at"`
}


// Token структура токена, требуется для подписи запросов клиента
type Token struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Platform  string     `json:"platform"`
	DeviceID  string     `json:"device_id"`
	Value     string     `json:"value"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// Note структура заметки
type Note struct {
	ID        int64      `json:"id"`
	UserID    int        `json:"user_id"`
	Text      string     `json:"text"`
	FileID    string     `json:"file_id"`
	FileType  byte       `json:"file_type"`
	CreatedAt *time.Time `json:"created_at"`
}
