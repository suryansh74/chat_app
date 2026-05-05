package chatrepositories

import (
	chatdomain "github.com/suryansh74/chat_app/internal/chat/domain"
)

type ChatRepositoryPort interface {
	CreateMessage(msg *chatdomain.Message) error
	GetMessages(userID, friendID string, limit, offset int) ([]chatdomain.MessageListItem, error)
	GetLastMessage(userID, friendID string) (*chatdomain.Message, error)
	SearchMessages(userID, query string) ([]chatdomain.MessageListItem, error)
	SearchConversationMessages(userID, friendID, query string) ([]chatdomain.MessageListItem, error)
	MarkAsRead(userID, friendID string) error
}
