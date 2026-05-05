package notificationrepositories

import (
	notificationdomain "github.com/suryansh74/chat_app/internal/notification/domain"
)

type NotificationRepositoryPort interface {
	Create(notification *notificationdomain.Notification) error
	GetByID(id string) (*notificationdomain.Notification, error)
	GetByUserID(userID string, limit, offset int) ([]notificationdomain.NotificationListItem, error)
	MarkAsRead(id string) error
	MarkAllAsRead(userID string) error
	Delete(id string) error
	GetUnreadCount(userID string) (int64, error)
}
