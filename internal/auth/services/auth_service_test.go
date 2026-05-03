package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	authdomain "github.com/suryansh74/chat_app/internal/auth/domain"
	"golang.org/x/crypto/bcrypt"
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

func (m *MockUserRepository) GetUserByEmail(email string) (*authdomain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authdomain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *authdomain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestValidateRegisterInput_NameTooShort(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:                 "ab",
		Email:                "test@example.com",
		Password:             "Password1!",
		PasswordConfirmation: "Password1!",
	}

	errors := service.ValidateRegisterInput(input)

	assert.NotEmpty(t, errors, "should return error for name less than 3 chars")
}

func TestValidateRegisterInput_InvalidEmail(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:                 "John",
		Email:                "invalid-email",
		Password:             "Password1!",
		PasswordConfirmation: "Password1!",
	}

	errors := service.ValidateRegisterInput(input)

	assert.NotEmpty(t, errors, "should return error for invalid email")
}

func TestValidateRegisterInput_PasswordMissingUppercase(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:                 "John",
		Email:                "test@example.com",
		Password:             "password1!",
		PasswordConfirmation: "password1!",
	}

	errors := service.ValidateRegisterInput(input)

	assert.NotEmpty(t, errors, "should return error for password missing uppercase")
}

func TestValidateRegisterInput_PasswordMissingNumber(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:                 "John",
		Email:                "test@example.com",
		Password:             "Password!",
		PasswordConfirmation: "Password!",
	}

	errors := service.ValidateRegisterInput(input)

	assert.NotEmpty(t, errors, "should return error for password missing number")
}

func TestValidateRegisterInput_PasswordMissingSpecialChar(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:                 "John",
		Email:                "test@example.com",
		Password:             "Password1",
		PasswordConfirmation: "Password1",
	}

	errors := service.ValidateRegisterInput(input)

	assert.NotEmpty(t, errors, "should return error for password missing special char")
}

func TestValidateRegisterInput_ValidInput(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:                 "John",
		Email:                "test@example.com",
		Password:             "Password1!",
		PasswordConfirmation: "Password1!",
	}

	errors := service.ValidateRegisterInput(input)

	assert.Empty(t, errors, "should not return errors for valid input")
}

func TestValidateRegisterInput_EmptyFields(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.RegisterInput{
		Name:                 "",
		Email:                "",
		Password:             "",
		PasswordConfirmation: "",
	}

	errors := service.ValidateRegisterInput(input)

	assert.NotEmpty(t, errors, "should return errors for empty required fields")
}

func TestLogin_Success(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password1!"), bcrypt.DefaultCost)
	user := &authdomain.User{
		ID:       "user-id-123",
		Name:     "John",
		Email:    "john@example.com",
		Password: string(hashedPassword),
	}
	repo.On("GetUserByEmail", "john@example.com").Return(user, nil)

	input := &authdomain.LoginInput{
		Email:    "john@example.com",
		Password: "Password1!",
	}

	resultUser, err := service.Login(input)

	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.Equal(t, "user-id-123", resultUser.ID)
	assert.Equal(t, "john@example.com", resultUser.Email)
}

func TestLogin_InvalidEmail(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	input := &authdomain.LoginInput{
		Email:    "invalid-email",
		Password: "Password1!",
	}

	errors := service.ValidateLoginInput(input)

	assert.NotEmpty(t, errors, "should return error for invalid email")
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	repo.On("GetUserByEmail", "notfound@example.com").Return(nil, assert.AnError)

	input := &authdomain.LoginInput{
		Email:    "notfound@example.com",
		Password: "Password1!",
	}

	_, err := service.Login(input)

	assert.Error(t, err)
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewAuthService(repo)

	hashedPassword := "$2a$10$YourHashedPasswordHere"
	user := &authdomain.User{
		ID:       "user-id-123",
		Name:     "John",
		Email:    "john@example.com",
		Password: hashedPassword,
	}
	repo.On("GetUserByEmail", "john@example.com").Return(user, nil)

	input := &authdomain.LoginInput{
		Email:    "john@example.com",
		Password: "WrongPassword1!",
	}

	_, err := service.Login(input)

	assert.Error(t, err)
}
