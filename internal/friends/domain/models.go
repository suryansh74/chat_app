package friendsdomain

import (
	"time"
)

type Friend struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	FriendID      string    `json:"friend_id"`
	FriendName    string    `json:"friend_name,omitempty"`
	FriendEmail   string    `json:"friend_email,omitempty"`
	LastMessageAt time.Time `json:"last_message_at"`
	UnreadCount   int       `json:"unread_count"`
	CreatedAt     time.Time `json:"created_at"`
}

type FriendRequest struct {
	ID         string    `json:"id"`
	FromUserID string    `json:"from_user_id"`
	ToUserID   string    `json:"to_user_id"`
	ToUserName string    `json:"to_user_name,omitempty"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type FriendListItem struct {
	FriendID      string    `json:"friend_id"`
	FriendName    string    `json:"friend_name"`
	FriendEmail   string    `json:"friend_email"`
	LastMessage   string    `json:"last_message"`
	LastMessageAt time.Time `json:"last_message_at"`
	UnreadCount   int       `json:"unread_count"`
}

type SendFriendRequestInput struct {
	ToUserID string `json:"to_user_id" validate:"required"`
}

type AcceptFriendRequestInput struct {
	RequestID string `json:"request_id" validate:"required"`
}
