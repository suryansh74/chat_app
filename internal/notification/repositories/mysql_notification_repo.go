package notificationrepositories

import (
	"fmt"
	"time"

	notificationdomain "github.com/suryansh74/chat_app/internal/notification/domain"
	"gorm.io/gorm"
)

type MySQLNotificationRepository struct {
	db *gorm.DB
}

type NotificationModel struct {
	ID         string `gorm:"primaryKey;size:36"`
	FromUserID string `gorm:"size:36;not null;index"`
	ToUserID   string `gorm:"size:36;not null;index"`
	Type       string `gorm:"size:50;not null"`
	Content    string `gorm:"type:text"`
	IsRead     bool   `gorm:"default:false"`
	CreatedAt  time.Time
}

func (NotificationModel) TableName() string {
	return "notifications"
}

func NewMySQLNotificationRepository(db *gorm.DB) *MySQLNotificationRepository {
	return &MySQLNotificationRepository{db: db}
}

func (r *MySQLNotificationRepository) Create(notification *notificationdomain.Notification) error {
	model := &NotificationModel{
		ID:         notification.ID,
		FromUserID: notification.FromUserID,
		ToUserID:   notification.ToUserID,
		Type:       notification.Type,
		Content:    notification.Content,
		IsRead:     notification.IsRead,
		CreatedAt:  notification.CreatedAt,
	}

	return r.db.Create(model).Error
}

func (r *MySQLNotificationRepository) GetByID(id string) (*notificationdomain.Notification, error) {
	var model NotificationModel
	err := r.db.Where("id = ?", id).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("notification not found")
		}
		return nil, err
	}

	return &notificationdomain.Notification{
		ID:         model.ID,
		FromUserID: model.FromUserID,
		ToUserID:   model.ToUserID,
		Type:       model.Type,
		Content:    model.Content,
		IsRead:     model.IsRead,
		CreatedAt:  model.CreatedAt,
	}, nil
}

func (r *MySQLNotificationRepository) GetByUserID(userID string, limit, offset int) ([]notificationdomain.NotificationListItem, error) {
	var results []struct {
		ID        string    `gorm:"column:id"`
		Type      string    `gorm:"column:type"`
		Content   string    `gorm:"column:content"`
		IsRead    bool      `gorm:"column:is_read"`
		FromUser  string    `gorm:"column:from_user"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}

	err := r.db.Table("notifications").
		Select("notifications.id, notifications.type, notifications.content, notifications.is_read, users.name as from_user, notifications.created_at").
		Joins("JOIN users ON users.id = notifications.from_user_id").
		Where("notifications.to_user_id = ?", userID).
		Order("notifications.created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	notifications := make([]notificationdomain.NotificationListItem, len(results))
	for i, res := range results {
		notifications[i] = notificationdomain.NotificationListItem{
			ID:        res.ID,
			Type:      res.Type,
			Content:   res.Content,
			IsRead:    res.IsRead,
			FromUser:  res.FromUser,
			CreatedAt: res.CreatedAt,
		}
	}

	return notifications, nil
}

func (r *MySQLNotificationRepository) MarkAsRead(id string) error {
	return r.db.Model(&NotificationModel{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *MySQLNotificationRepository) MarkAllAsRead(userID string) error {
	return r.db.Model(&NotificationModel{}).Where("to_user_id = ?", userID).Update("is_read", true).Error
}

func (r *MySQLNotificationRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&NotificationModel{}).Error
}

func (r *MySQLNotificationRepository) GetUnreadCount(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&NotificationModel{}).Where("to_user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&NotificationModel{})
}
