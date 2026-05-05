package chatrepositories

import (
	"fmt"
	"time"

	chatdomain "github.com/suryansh74/chat_app/internal/chat/domain"
	"gorm.io/gorm"
)

type MySQLChatRepository struct {
	db *gorm.DB
}

type MessageModel struct {
	ID         string    `gorm:"primaryKey;size:36"`
	FromUserID string    `gorm:"size:36;not null;index"`
	ToUserID   string    `gorm:"size:36;not null;index"`
	Content    string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"index"`
}

func (MessageModel) TableName() string {
	return "messages"
}

func NewMySQLChatRepository(db *gorm.DB) *MySQLChatRepository {
	return &MySQLChatRepository{db: db}
}

func (r *MySQLChatRepository) CreateMessage(msg *chatdomain.Message) error {
	model := &MessageModel{
		ID:         msg.ID,
		FromUserID: msg.FromUserID,
		ToUserID:   msg.ToUserID,
		Content:    msg.Content,
		CreatedAt:  msg.CreatedAt,
	}

	return r.db.Create(model).Error
}

func (r *MySQLChatRepository) GetMessages(userID, friendID string, limit, offset int) ([]chatdomain.MessageListItem, error) {
	var results []struct {
		ID         string    `gorm:"column:id"`
		FromUserID string    `gorm:"column:from_user_id"`
		ToUserID   string    `gorm:"column:to_user_id"`
		Content    string    `gorm:"column:content"`
		CreatedAt  time.Time `gorm:"column:created_at"`
	}

	err := r.db.Table("messages").
		Select("id, from_user_id, to_user_id, content, created_at").
		Where("((from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?))", userID, friendID, friendID, userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	messages := make([]chatdomain.MessageListItem, len(results))
	for i, res := range results {
		messages[i] = chatdomain.MessageListItem{
			MessageID:  res.ID,
			FromUserID: res.FromUserID,
			ToUserID:   res.ToUserID,
			Content:    res.Content,
			IsMe:       res.FromUserID == userID,
			CreatedAt:  res.CreatedAt,
		}
	}

	return messages, nil
}

func (r *MySQLChatRepository) GetLastMessage(userID, friendID string) (*chatdomain.Message, error) {
	var model MessageModel
	err := r.db.Where("((from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?))", userID, friendID, friendID, userID).
		Order("created_at DESC").
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no messages")
		}
		return nil, err
	}

	return &chatdomain.Message{
		ID:         model.ID,
		FromUserID: model.FromUserID,
		ToUserID:   model.ToUserID,
		Content:    model.Content,
		CreatedAt:  model.CreatedAt,
	}, nil
}

func (r *MySQLChatRepository) SearchMessages(userID, query string) ([]chatdomain.MessageListItem, error) {
	var results []struct {
		ID         string    `gorm:"column:id"`
		FromUserID string    `gorm:"column:from_user_id"`
		ToUserID   string    `gorm:"column:to_user_id"`
		Content    string    `gorm:"column:content"`
		CreatedAt  time.Time `gorm:"column:created_at"`
	}

	err := r.db.Table("messages").
		Select("id, from_user_id, to_user_id, content, created_at").
		Where("((from_user_id = ? OR to_user_id = ?) AND content LIKE ?)", userID, userID, "%"+query+"%").
		Order("created_at DESC").
		Limit(50).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	messages := make([]chatdomain.MessageListItem, len(results))
	for i, res := range results {
		messages[i] = chatdomain.MessageListItem{
			MessageID:  res.ID,
			FromUserID: res.FromUserID,
			ToUserID:   res.ToUserID,
			Content:    res.Content,
			IsMe:       res.FromUserID == userID,
			CreatedAt:  res.CreatedAt,
		}
	}

	return messages, nil
}

func (r *MySQLChatRepository) SearchConversationMessages(userID, friendID, query string) ([]chatdomain.MessageListItem, error) {
	var results []struct {
		ID         string    `gorm:"column:id"`
		FromUserID string    `gorm:"column:from_user_id"`
		ToUserID   string    `gorm:"column:to_user_id"`
		Content    string    `gorm:"column:content"`
		CreatedAt  time.Time `gorm:"column:created_at"`
	}

	err := r.db.Table("messages").
		Select("id, from_user_id, to_user_id, content, created_at").
		Where("((from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)) AND content LIKE ?", userID, friendID, friendID, userID, "%"+query+"%").
		Order("created_at DESC").
		Limit(50).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	messages := make([]chatdomain.MessageListItem, len(results))
	for i, res := range results {
		messages[i] = chatdomain.MessageListItem{
			MessageID:  res.ID,
			FromUserID: res.FromUserID,
			ToUserID:   res.ToUserID,
			Content:    res.Content,
			IsMe:       res.FromUserID == userID,
			CreatedAt:  res.CreatedAt,
		}
	}

	return messages, nil
}

func (r *MySQLChatRepository) MarkAsRead(userID, friendID string) error {
	return r.db.Model(&MessageModel{}).Where("to_user_id = ? AND from_user_id = ?", userID, friendID).Update("is_read", true).Error
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&MessageModel{})
}
