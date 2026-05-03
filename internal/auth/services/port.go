package auth

import authdomain "github.com/suryansh74/chat_app/internal/auth/domain"

type AuthServicePort interface {
	Register(input *authdomain.RegisterInput) (*authdomain.User, error)
	ValidateRegisterInput(input *authdomain.RegisterInput) []error
	Login(input *authdomain.LoginInput) (*authdomain.User, error)
	ValidateLoginInput(input *authdomain.LoginInput) []error
	Logout() error
}
