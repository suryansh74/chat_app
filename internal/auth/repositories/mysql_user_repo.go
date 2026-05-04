package repositories

import (
	"fmt"
	"time"

	authdomain "github.com/suryansh74/chat_app/internal/auth/domain"
	"gorm.io/gorm"
)

type MySQLUserRepository struct {
	db *gorm.DB
}

type UserModel struct {
	ID                   string `gorm:"primaryKey;size:36"`
	Name                 string `gorm:"size:255;not null"`
	Email                string `gorm:"size:255;not null;uniqueIndex"`
	Password             string `gorm:"size:255;not null"`
	IsVerified           bool   `gorm:"default:false"`
	OTP                  string `gorm:"size:6"`
	OTPExpiry            *time.Time
	OTPAttempts          int    `gorm:"default:0"`
	PasswordResetOTP     string `gorm:"size:6"`
	PasswordResetExpiry  *time.Time
	PasswordResetAttempt int `gorm:"default:0"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (UserModel) TableName() string {
	return "users"
}

func NewMySQLUserRepository(db *gorm.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&UserModel{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *MySQLUserRepository) CreateUser(user *authdomain.User) error {
	userModel := &UserModel{
		ID:                   user.ID,
		Name:                 user.Name,
		Email:                user.Email,
		Password:             user.Password,
		IsVerified:           user.IsVerified,
		OTP:                  user.OTP,
		OTPExpiry:            &user.OTPExpiry,
		OTPAttempts:          user.OTPAttempts,
		PasswordResetOTP:     user.PasswordResetOTP,
		PasswordResetExpiry:  &user.PasswordResetExpiry,
		PasswordResetAttempt: user.PasswordResetAttempt,
	}

	return r.db.Create(userModel).Error
}

func (r *MySQLUserRepository) GetUserByEmail(email string) (*authdomain.User, error) {
	var userModel UserModel
	err := r.db.Where("email = ?", email).First(&userModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	user := &authdomain.User{
		ID:                   userModel.ID,
		Name:                 userModel.Name,
		Email:                userModel.Email,
		Password:             userModel.Password,
		IsVerified:           userModel.IsVerified,
		OTP:                  userModel.OTP,
		OTPAttempts:          userModel.OTPAttempts,
		PasswordResetOTP:     userModel.PasswordResetOTP,
		PasswordResetAttempt: userModel.PasswordResetAttempt,
	}

	if userModel.OTPExpiry != nil {
		user.OTPExpiry = *userModel.OTPExpiry
	}
	if userModel.PasswordResetExpiry != nil {
		user.PasswordResetExpiry = *userModel.PasswordResetExpiry
	}

	return user, nil
}

func (r *MySQLUserRepository) UpdateUser(user *authdomain.User) error {
	userModel := &UserModel{
		ID:                   user.ID,
		Name:                 user.Name,
		Email:                user.Email,
		Password:             user.Password,
		IsVerified:           user.IsVerified,
		OTP:                  user.OTP,
		OTPExpiry:            &user.OTPExpiry,
		OTPAttempts:          user.OTPAttempts,
		PasswordResetOTP:     user.PasswordResetOTP,
		PasswordResetExpiry:  &user.PasswordResetExpiry,
		PasswordResetAttempt: user.PasswordResetAttempt,
	}

	return r.db.Save(userModel).Error
}

// AutoMigrate creates the users table if it doesn't exist
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&UserModel{})
}

// GetDSN returns the MySQL connection string
func GetDSN(host, port, user, password, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)
}
