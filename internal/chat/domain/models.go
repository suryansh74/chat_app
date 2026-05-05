package chatdomain

import (
	"time"
)

type Message struct {
	ID         string    `json:"id"`
	FromUserID string    `json:"from_user_id"`
	ToUserID   string    `json:"to_user_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type MessageListItem struct {
	MessageID  string    `json:"message_id"`
	FromUserID string    `json:"from_user_id"`
	ToUserID   string    `json:"to_user_id"`
	Content    string    `json:"content"`
	IsMe       bool      `json:"is_me"`
	CreatedAt  time.Time `json:"created_at"`
}

type SendMessageInput struct {
	ToUserID string `json:"to_user_id" validate:"required"`
	Content  string `json:"content" validate:"required"`
}
