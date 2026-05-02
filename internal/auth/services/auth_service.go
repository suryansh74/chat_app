package auth

import (
	"github.com/go-playground/validator/v10"
	authdomain "github.com/suryansh74/chat_app/internal/auth/domain"
	"github.com/suryansh74/chat_app/internal/auth/repositories"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Message
}

type AuthService struct {
	repo     repositories.UserRepository
	validate *validator.Validate
}

func NewAuthService(repo repositories.UserRepository) *AuthService {
	return &AuthService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *AuthService) Register(input *authdomain.RegisterInput) error {
	errors := s.ValidateRegisterInput(input)
	if len(errors) > 0 {
		return errors[0]
	}

	exists, err := s.repo.EmailExists(input.Email)
	if err != nil {
		return err
	}
	if exists {
		return ValidationError{Field: "Email", Message: "email already exists"}
	}

	user := &authdomain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ValidateRegisterInput(input *authdomain.RegisterInput) []ValidationError {
	var errors []ValidationError

	err := s.validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field:   err.Field(),
				Message: err.Error(),
			})
		}
	}

	if len(input.Password) > 0 {
		hasUpper, hasDigit, hasSpecial := false, false, false
		for _, ch := range input.Password {
			switch {
			case ch >= 'A' && ch <= 'Z':
				hasUpper = true
			case ch >= '0' && ch <= '9':
				hasDigit = true
			case ch == '!' || ch == '@' || ch == '#' || ch == '$' || ch == '%' || ch == '^' || ch == '&' || ch == '*':
				hasSpecial = true
			}
		}
		if !(hasUpper && hasDigit && hasSpecial) {
			errors = append(errors, ValidationError{
				Field:   "Password",
				Message: "password must contain at least one uppercase letter, one number, and one special character",
			})
		}
	}

	return errors
}
