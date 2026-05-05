package friendsrepositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	friendsdomain "github.com/suryansh74/chat_app/internal/friends/domain"
	"gorm.io/gorm"
)

type MySQLFriendRepository struct {
	db *gorm.DB
}

type FriendModel struct {
	ID            string    `gorm:"primaryKey;size:36"`
	UserID        string    `gorm:"size:36;not null;index:idx_user_friend"`
	FriendID      string    `gorm:"size:36;not null;index:idx_user_friend"`
	LastMessageAt time.Time `gorm:"index"`
	UnreadCount   int       `gorm:"default:0"`
	CreatedAt     time.Time
}

type UserModel struct {
	ID    string `gorm:"primaryKey;size:36"`
	Name  string `gorm:"size:255;not null"`
	Email string `gorm:"size:255;not null;uniqueIndex"`
}

func (UserModel) TableName() string {
	return "users"
}

func (FriendModel) TableName() string {
	return "friends"
}

func NewMySQLFriendRepository(db *gorm.DB) *MySQLFriendRepository {
	return &MySQLFriendRepository{db: db}
}

func (r *MySQLFriendRepository) CreateFriend(userID, friendID string) error {
	friend := &FriendModel{
		ID:       uuid.New().String(),
		UserID:   userID,
		FriendID: friendID,
	}

	return r.db.Create(friend).Error
}

func (r *MySQLFriendRepository) DeleteFriend(userID, friendID string) error {
	return r.db.Where("user_id = ? AND friend_id = ?", userID, friendID).Delete(&FriendModel{}).Error
}

func (r *MySQLFriendRepository) GetFriendsByUserID(userID string) ([]friendsdomain.FriendListItem, error) {
	var results []struct {
		FriendID      string    `gorm:"column:friend_id"`
		FriendName    string    `gorm:"column:name"`
		FriendEmail   string    `gorm:"column:email"`
		LastMessageAt time.Time `gorm:"column:last_message_at"`
		UnreadCount   int       `gorm:"column:unread_count"`
		CreatedAt     time.Time
	}

	err := r.db.Table("friends").
		Select("friends.friend_id, users.name, users.email, last_message_at, unread_count, friends.created_at").
		Joins("JOIN users ON users.id = friends.friend_id").
		Where("friends.user_id = ?", userID).
		Order("last_message_at DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	friends := make([]friendsdomain.FriendListItem, len(results))
	for i, res := range results {
		friends[i] = friendsdomain.FriendListItem{
			FriendID:      res.FriendID,
			FriendName:    res.FriendName,
			FriendEmail:   res.FriendEmail,
			LastMessageAt: res.LastMessageAt,
			UnreadCount:   res.UnreadCount,
		}
	}

	return friends, nil
}

func (r *MySQLFriendRepository) GetFriendByID(userID, friendID string) (*friendsdomain.Friend, error) {
	var friend FriendModel
	err := r.db.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&friend).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("friend not found")
		}
		return nil, err
	}

	return &friendsdomain.Friend{
		ID:            friend.ID,
		UserID:        friend.UserID,
		FriendID:      friend.FriendID,
		LastMessageAt: friend.LastMessageAt,
		UnreadCount:   friend.UnreadCount,
		CreatedAt:     friend.CreatedAt,
	}, nil
}

func (r *MySQLFriendRepository) IsFriend(userID, friendID string) (bool, error) {
	var count int64
	err := r.db.Model(&FriendModel{}).Where("user_id = ? AND friend_id = ?", userID, friendID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *MySQLFriendRepository) GetFriendsByEmail(email string) ([]friendsdomain.Friend, error) {
	var results []struct {
		ID       string `gorm:"column:id"`
		FriendID string `gorm:"column:friend_id"`
		Name     string `gorm:"column:name"`
		Email    string `gorm:"column:email"`
	}

	// Use prefix matching for better UX - search for emails starting with query
	// Also include exact match for better results
	err := r.db.Table("users").
		Select("users.id, users.id as friend_id, users.name, users.email").
		Where("email LIKE ? OR name LIKE ?", email+"%", email+"%").
		Order("email ASC").
		Limit(5).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// If no prefix matches, try contains
	if len(results) == 0 {
		err = r.db.Table("users").
			Select("users.id, users.id as friend_id, users.name, users.email").
			Where("email LIKE ? OR name LIKE ?", "%"+email+"%", "%"+email+"%").
			Order("email ASC").
			Limit(5).
			Scan(&results).Error

		if err != nil {
			return nil, err
		}
	}

	friends := make([]friendsdomain.Friend, len(results))
	for i, res := range results {
		friends[i] = friendsdomain.Friend{
			ID:          res.ID,
			FriendID:    res.FriendID,
			FriendName:  res.Name,
			FriendEmail: res.Email,
		}
	}

	return friends, nil
}

func (r *MySQLFriendRepository) GetFriendByEmail(email string) (*friendsdomain.Friend, error) {
	var user UserModel
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &friendsdomain.Friend{
		ID:          user.ID,
		FriendID:    user.ID,
		FriendName:  user.Name,
		FriendEmail: user.Email,
	}, nil
}

func (r *MySQLFriendRepository) GetUserInfo(userID string) (*UserInfo, error) {
	var user UserModel
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &UserInfo{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&FriendModel{})
}
