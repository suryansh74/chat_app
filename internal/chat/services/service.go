package chatservices

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	chatdomain "github.com/suryansh74/chat_app/internal/chat/domain"
	chatrepositories "github.com/suryansh74/chat_app/internal/chat/repositories"
	friendsrepositories "github.com/suryansh74/chat_app/internal/friends/repositories"
	notificationservices "github.com/suryansh74/chat_app/internal/notification/services"
	ws "github.com/suryansh74/chat_app/internal/ws"
)

type ChatServicePort interface {
	SendMessage(fromUserID, toUserID, content string) (*chatdomain.Message, error)
	GetMessages(userID, friendID string, limit, offset int) ([]chatdomain.MessageListItem, error)
	SearchMessages(userID, query string) ([]chatdomain.MessageListItem, error)
	SearchConversationMessages(userID, friendID, query string) ([]chatdomain.MessageListItem, error)
}

type chatService struct {
	chatRepo        chatrepositories.ChatRepositoryPort
	friendRepo      friendsrepositories.FriendRepositoryPort
	notificationSvc notificationservices.NotificationServicePort
	wsHub           *ws.Hub
}

func NewChatService(
	chatRepo chatrepositories.ChatRepositoryPort,
	friendRepo friendsrepositories.FriendRepositoryPort,
	notificationSvc notificationservices.NotificationServicePort,
	wsHub *ws.Hub,
) ChatServicePort {
	return &chatService{
		chatRepo:        chatRepo,
		friendRepo:      friendRepo,
		notificationSvc: notificationSvc,
		wsHub:           wsHub,
	}
}

func (s *chatService) SendMessage(fromUserID, toUserID, content string) (*chatdomain.Message, error) {
	isFriend, err := s.friendRepo.IsFriend(fromUserID, toUserID)
	if err != nil {
		return nil, err
	}
	if !isFriend {
		return nil, fmt.Errorf("not friends")
	}

	msg := &chatdomain.Message{
		ID:         uuid.New().String(),
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Content:    content,
		CreatedAt:  time.Now(),
	}

	err = s.chatRepo.CreateMessage(msg)
	if err != nil {
		return nil, err
	}

	if s.wsHub != nil {
		s.wsHub.SendToUserJSON(toUserID, "new_message", msg)
	}

	return msg, nil
}

func (s *chatService) GetMessages(userID, friendID string, limit, offset int) ([]chatdomain.MessageListItem, error) {
	return s.chatRepo.GetMessages(userID, friendID, limit, offset)
}

func (s *chatService) SearchMessages(userID, query string) ([]chatdomain.MessageListItem, error) {
	return s.chatRepo.SearchMessages(userID, query)
}

func (s *chatService) SearchConversationMessages(userID, friendID, query string) ([]chatdomain.MessageListItem, error) {
	return s.chatRepo.SearchConversationMessages(userID, friendID, query)
}
