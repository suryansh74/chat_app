package notificationservices

import (
	"fmt"

	"github.com/google/uuid"
	notificationdomain "github.com/suryansh74/chat_app/internal/notification/domain"
	notificationrepositories "github.com/suryansh74/chat_app/internal/notification/repositories"
	"github.com/suryansh74/chat_app/shared/cache"
	emailadapters "github.com/suryansh74/chat_app/shared/email/adapters"
	"time"
)

type NotificationServicePort interface {
	CreateFriendRequest(fromUserID, toUserID string) error
	CreateFriendAccepted(fromUserID, toUserID string) error
	CreateFriendRejected(fromUserID, toUserID string) error
	GetNotifications(userID string, limit, offset int) ([]notificationdomain.NotificationListItem, error)
	GetNotificationByID(id string) (*notificationdomain.Notification, error)
	MarkAsRead(id string) error
	MarkAllAsRead(userID string) error
	DeleteNotification(id string) error
	GetUnreadCount(userID string) (int64, error)
}

type notificationService struct {
	repo        notificationrepositories.NotificationRepositoryPort
	cache       cache.CachePort
	emailSender emailadapters.EmailSenderPort
}

func NewNotificationService(
	repo notificationrepositories.NotificationRepositoryPort,
	cache cache.CachePort,
	emailSender emailadapters.EmailSenderPort,
) NotificationServicePort {
	return &notificationService{
		repo:        repo,
		cache:       cache,
		emailSender: emailSender,
	}
}

func unreadCountKey(userID string) string {
	return fmt.Sprintf("unread_count:%s", userID)
}

func (s *notificationService) CreateFriendRequest(fromUserID, toUserID string) error {
	notification := &notificationdomain.Notification{
		ID:         uuid.New().String(),
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Type:       notificationdomain.NotificationTypeFriendRequest,
		Content:    "sent you a friend request",
		IsRead:     false,
		CreatedAt:  time.Now(),
	}

	err := s.repo.Create(notification)
	if err != nil {
		return err
	}

	s.cache.Delete(unreadCountKey(toUserID))
	return nil
}

func (s *notificationService) CreateFriendAccepted(fromUserID, toUserID string) error {
	notification := &notificationdomain.Notification{
		ID:         uuid.New().String(),
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Type:       notificationdomain.NotificationTypeFriendAccepted,
		Content:    "accepted your friend request",
		IsRead:     false,
		CreatedAt:  time.Now(),
	}

	return s.repo.Create(notification)
}

func (s *notificationService) CreateFriendRejected(fromUserID, toUserID string) error {
	notification := &notificationdomain.Notification{
		ID:         uuid.New().String(),
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Type:       "FRIEND_REJECTED",
		Content:    "rejected your friend request",
		IsRead:     false,
		CreatedAt:  time.Now(),
	}

	return s.repo.Create(notification)
}

func (s *notificationService) GetNotifications(userID string, limit, offset int) ([]notificationdomain.NotificationListItem, error) {
	return s.repo.GetByUserID(userID, limit, offset)
}

func (s *notificationService) GetNotificationByID(id string) (*notificationdomain.Notification, error) {
	return s.repo.GetByID(id)
}

func (s *notificationService) MarkAsRead(id string) error {
	err := s.repo.MarkAsRead(id)
	if err != nil {
		return err
	}

	notif, err := s.repo.GetByID(id)
	if err == nil && notif != nil {
		s.cache.Delete(unreadCountKey(notif.ToUserID))
	}

	return nil
}

func (s *notificationService) MarkAllAsRead(userID string) error {
	err := s.repo.MarkAllAsRead(userID)
	if err != nil {
		return err
	}

	s.cache.Delete(unreadCountKey(userID))
	return nil
}

func (s *notificationService) DeleteNotification(id string) error {
	notif, err := s.repo.GetByID(id)
	if err == nil && notif != nil {
		s.cache.Delete(unreadCountKey(notif.ToUserID))
	}

	return s.repo.Delete(id)
}

func (s *notificationService) GetUnreadCount(userID string) (int64, error) {
	cached, err := s.cache.Get(unreadCountKey(userID))
	if err == nil && cached != nil {
		if val, ok := cached.(float64); ok {
			return int64(val), nil
		}
	}

	count, err := s.repo.GetUnreadCount(userID)
	if err != nil {
		return 0, err
	}

	s.cache.Set(unreadCountKey(userID), count, 10*time.Second)
	return count, nil
}
