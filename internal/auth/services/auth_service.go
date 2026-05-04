package auth

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	apperr "github.com/suryansh74/chat_app/internal/auth/apperr"
	authdomain "github.com/suryansh74/chat_app/internal/auth/domain"
	"github.com/suryansh74/chat_app/internal/auth/repositories"
	"github.com/suryansh74/chat_app/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

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

func (s *AuthService) Register(input *authdomain.RegisterInput) (*authdomain.User, error) {
	errors := s.ValidateRegisterInput(input)
	if len(errors) > 0 {
		return nil, errors[0]
	}

	exists, err := s.repo.EmailExists(input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperr.NewEmailAlreadyExists(input.Email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &authdomain.User{
		ID:       uuid.New().String(),
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(input *authdomain.LoginInput) (*authdomain.User, error) {
	errors := s.ValidateLoginInput(input)
	if len(errors) > 0 {
		return nil, errors[0]
	}

	user, err := s.repo.GetUserByEmail(input.Email)
	if err != nil {
		logger.Log.Warn("Login: user not found", "email", input.Email, "error", err.Error())
		return nil, ErrInvalidCredentials
	}

	if user == nil {
		logger.Log.Warn("Login: user is nil after search", "email", input.Email)
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		logger.Log.Warn("Login: password mismatch", "email", input.Email)
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *AuthService) ValidateLoginInput(input *authdomain.LoginInput) []error {
	var validationErrors []ValidationError

	err := s.validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Message: err.Error(),
			})
		}
	}

	var errors []error
	for _, e := range validationErrors {
		errors = append(errors, e)
	}
	return errors
}

func (s *AuthService) ValidateRegisterInput(input *authdomain.RegisterInput) []error {
	var validationErrors []ValidationError

	err := s.validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
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
			validationErrors = append(validationErrors, ValidationError{
				Field:   "Password",
				Message: "password must contain at least one uppercase letter, one number, and one special character",
			})
		}
	}

	if input.Password != input.PasswordConfirmation {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "PasswordConfirmation",
			Message: "password and password_confirmation must match",
		})
	}

	var errors []error
	for _, e := range validationErrors {
		errors = append(errors, e)
	}
	return errors
}

func (s *AuthService) Logout() error {
	return nil
}

func (s *AuthService) ValidateSetPassword(input *authdomain.SetPasswordInput) []error {
	var validationErrors []ValidationError

	err := s.validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
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
			validationErrors = append(validationErrors, ValidationError{
				Field:   "Password",
				Message: "password must contain at least one uppercase letter, one number, and one special character",
			})
		}
	}

	if input.Password != input.PasswordConfirmation {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "PasswordConfirmation",
			Message: "password and password_confirmation must match",
		})
	}

	var errors []error
	for _, e := range validationErrors {
		errors = append(errors, e)
	}
	return errors
}
