package repositories

import authdomain "github.com/suryansh74/chat_app/internal/auth/domain"

type UserRepository interface {
	EmailExists(email string) (bool, error)
	CreateUser(user *authdomain.User) error
	GetUserByEmail(email string) (*authdomain.User, error)
}
