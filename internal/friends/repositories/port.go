package friendsrepositories

import (
	friendsdomain "github.com/suryansh74/chat_app/internal/friends/domain"
)

type FriendRepositoryPort interface {
	CreateFriend(userID, friendID string) error
	DeleteFriend(userID, friendID string) error
	GetFriendsByUserID(userID string) ([]friendsdomain.FriendListItem, error)
	GetFriendByID(userID, friendID string) (*friendsdomain.Friend, error)
	IsFriend(userID, friendID string) (bool, error)
	GetFriendsByEmail(email string) ([]friendsdomain.Friend, error)
	GetFriendByEmail(email string) (*friendsdomain.Friend, error)
	GetUserInfo(userID string) (*UserInfo, error)
}

type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
