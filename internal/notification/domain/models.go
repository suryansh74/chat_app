package notificationdomain

import (
	"time"
)

type Notification struct {
	ID         string    `json:"id"`
	FromUserID string    `json:"from_user_id"`
	ToUserID   string    `json:"to_user_id"`
	Type       string    `json:"type"`
	Content    string    `json:"content"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

type NotificationListItem struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	IsRead    bool      `json:"is_read"`
	FromUser  string    `json:"from_user"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateFriendRequestInput struct {
	ToUserID string `json:"to_user_id" validate:"required"`
}

type CreateFriendAcceptedInput struct {
	ToUserID string `json:"to_user_id" validate:"required"`
}
