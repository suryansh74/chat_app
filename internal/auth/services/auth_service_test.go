package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	authdomain "github.com/suryansh74/chat_app/internal/auth/domain"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) EmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *authdomain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestValidateRegisterInput_NameTooShort(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:     "ab",
		Email:    "test@example.com",
		Password: "Password1!",
	}

	errors := service.ValidateRegisterInput(input)

	assert.NotEmpty(t, errors, "should return error for name less than 3 chars")
}

func TestValidateRegisterInput_InvalidEmail(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:     "John",
		Email:    "invalid-email",
		Password: "Password1!",
	}

	errors := service.ValidateRegisterInput(input)

	assert.NotEmpty(t, errors, "should return error for invalid email")
}

func TestValidateRegisterInput_PasswordMissingUppercase(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:     "John",
		Email:    "test@example.com",
		Password: "password1!",
	}

	errors := service.ValidateRegisterInput(input)

	assert.NotEmpty(t, errors, "should return error for password missing uppercase")
}

func TestValidateRegisterInput_PasswordMissingNumber(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:     "John",
		Email:    "test@example.com",
		Password: "Password!",
	}

	errors := service.ValidateRegisterInput(input)

	assert.NotEmpty(t, errors, "should return error for password missing number")
}

func TestValidateRegisterInput_PasswordMissingSpecialChar(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:     "John",
		Email:    "test@example.com",
		Password: "Password1",
	}

	errors := service.ValidateRegisterInput(input)

	assert.NotEmpty(t, errors, "should return error for password missing special char")
}

func TestValidateRegisterInput_ValidInput(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:     "John",
		Email:    "test@example.com",
		Password: "Password1!",
	}

	errors := service.ValidateRegisterInput(input)

	assert.Empty(t, errors, "should not return errors for valid input")
}

func TestValidateRegisterInput_EmptyFields(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:     "",
		Email:    "",
		Password: "",
	}

	errors := service.ValidateRegisterInput(input)

	assert.NotEmpty(t, errors, "should return errors for empty required fields")
}
