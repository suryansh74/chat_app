package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authdomain "github.com/suryansh74/chat_app/internal/auth/domain"
)

type MockUserRepositoryForOTP struct {
	users map[string]*authdomain.User
}

func NewMockUserRepositoryForOTP() *MockUserRepositoryForOTP {
	return &MockUserRepositoryForOTP{
		users: make(map[string]*authdomain.User),
	}
}

func (r *MockUserRepositoryForOTP) EmailExists(email string) (bool, error) {
	_, exists := r.users[email]
	return exists, nil
}

func (r *MockUserRepositoryForOTP) CreateUser(user *authdomain.User) error {
	r.users[user.Email] = user
	return nil
}

func (r *MockUserRepositoryForOTP) GetUserByEmail(email string) (*authdomain.User, error) {
	user, exists := r.users[email]
	if !exists {
		return nil, assert.AnError
	}
	return user, nil
}

func (r *MockUserRepositoryForOTP) UpdateUser(user *authdomain.User) error {
	r.users[user.Email] = user
	return nil
}

func TestEmailVerificationService_GenerateOTP_Generates6Digits(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 5, 3)

	otp, err := service.GenerateOTP()

	require.NoError(t, err)
	assert.Len(t, otp, 6, "OTP should be 6 digits")

	for _, c := range otp {
		assert.GreaterOrEqual(t, c, '0')
		assert.LessOrEqual(t, c, '9')
	}
}

func TestEmailVerificationService_SendOTP_StoresOTPInUser(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 5, 3)

	user := &authdomain.User{
		ID:         uuid.New().String(),
		Name:       "John",
		Email:      "john@example.com",
		IsVerified: false,
	}
	repo.CreateUser(user)

	otp, err := service.SendOTP("john@example.com")

	require.NoError(t, err)
	assert.NotEmpty(t, otp)

	updatedUser, _ := repo.GetUserByEmail("john@example.com")
	assert.Equal(t, otp, updatedUser.OTP)
	assert.NotZero(t, updatedUser.OTPExpiry)
}

func TestEmailVerificationService_VerifyOTP_Success(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 5, 3)

	user := &authdomain.User{
		ID:         uuid.New().String(),
		Name:       "John",
		Email:      "john@example.com",
		IsVerified: false,
	}
	repo.CreateUser(user)

	otp, _ := service.SendOTP("john@example.com")

	err := service.VerifyOTP("john@example.com", otp)

	require.NoError(t, err)

	verifiedUser, _ := repo.GetUserByEmail("john@example.com")
	assert.True(t, verifiedUser.IsVerified)
	assert.Empty(t, verifiedUser.OTP)
}

func TestEmailVerificationService_VerifyOTP_WrongOTP(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 5, 3)

	user := &authdomain.User{
		ID:         uuid.New().String(),
		Name:       "John",
		Email:      "john@example.com",
		IsVerified: false,
	}
	repo.CreateUser(user)

	service.SendOTP("john@example.com")

	err := service.VerifyOTP("john@example.com", "000000")

	assert.Equal(t, ErrInvalidOTP, err)

	verifiedUser, _ := repo.GetUserByEmail("john@example.com")
	assert.False(t, verifiedUser.IsVerified)
}

func TestEmailVerificationService_VerifyOTP_IncrementsAttempts(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 5, 3)

	user := &authdomain.User{
		ID:         uuid.New().String(),
		Name:       "John",
		Email:      "john@example.com",
		IsVerified: false,
	}
	repo.CreateUser(user)

	service.SendOTP("john@example.com")

	service.VerifyOTP("john@example.com", "000000")
	service.VerifyOTP("john@example.com", "000000")

	verifiedUser, _ := repo.GetUserByEmail("john@example.com")
	assert.Equal(t, 2, verifiedUser.OTPAttempts)
}

func TestEmailVerificationService_VerifyOTP_MaxAttemptsBlocks(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 5, 3)

	user := &authdomain.User{
		ID:         uuid.New().String(),
		Name:       "John",
		Email:      "john@example.com",
		IsVerified: false,
	}
	repo.CreateUser(user)

	service.SendOTP("john@example.com")

	err1 := service.VerifyOTP("john@example.com", "000000")
	err2 := service.VerifyOTP("john@example.com", "000000")
	err3 := service.VerifyOTP("john@example.com", "000000")

	assert.Equal(t, ErrInvalidOTP, err1)
	assert.Equal(t, ErrInvalidOTP, err2)
	assert.Equal(t, ErrMaxAttempts, err3)

	verifiedUser, _ := repo.GetUserByEmail("john@example.com")
	assert.Equal(t, 0, verifiedUser.OTPAttempts)
	assert.Empty(t, verifiedUser.OTP)
}

func TestEmailVerificationService_VerifyOTP_ExpiredOTP(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 0, 3)

	user := &authdomain.User{
		ID:         uuid.New().String(),
		Name:       "John",
		Email:      "john@example.com",
		IsVerified: false,
	}
	repo.CreateUser(user)

	_, _ = service.SendOTP("john@example.com")

	time.Sleep(10 * time.Millisecond)

	err := service.VerifyOTP("john@example.com", "invalid")

	assert.Equal(t, ErrOTPExpired, err)
}

func TestEmailVerificationService_VerifyOTP_AlreadyVerified(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 5, 3)

	user := &authdomain.User{
		ID:         uuid.New().String(),
		Name:       "John",
		Email:      "john@example.com",
		IsVerified: true,
	}
	repo.CreateUser(user)

	_, err := service.SendOTP("john@example.com")

	assert.Equal(t, ErrAlreadyVerified, err)
}

func TestEmailVerificationService_PasswordResetOTP(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 5, 3)

	user := &authdomain.User{
		ID:    uuid.New().String(),
		Name:  "John",
		Email: "john@example.com",
	}
	repo.CreateUser(user)

	otp, err := service.SendPasswordResetOTP("john@example.com")

	require.NoError(t, err)
	assert.NotEmpty(t, otp)

	updatedUser, _ := repo.GetUserByEmail("john@example.com")
	assert.Equal(t, otp, updatedUser.PasswordResetOTP)
}

func TestEmailVerificationService_SetPassword(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 5, 3)

	user := &authdomain.User{
		ID:    uuid.New().String(),
		Name:  "John",
		Email: "john@example.com",
	}
	repo.CreateUser(user)

	err := service.SetPassword("john@example.com", "NewPassword1!")

	require.NoError(t, err)

	updatedUser, _ := repo.GetUserByEmail("john@example.com")
	assert.NotEqual(t, "NewPassword1!", updatedUser.Password)
}

func TestEmailVerificationService_SetPassword_EmptyPassword(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 5, 3)

	user := &authdomain.User{
		ID:    uuid.New().String(),
		Name:  "John",
		Email: "john@example.com",
	}
	repo.CreateUser(user)

	err := service.SetPassword("john@example.com", "")

	assert.NoError(t, err)
}

func TestEmailVerificationService_VerifyPasswordResetOTP_Success(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 5, 3)

	user := &authdomain.User{
		ID:    uuid.New().String(),
		Name:  "John",
		Email: "john@example.com",
	}
	repo.CreateUser(user)

	otp, _ := service.SendPasswordResetOTP("john@example.com")

	err := service.VerifyPasswordResetOTP("john@example.com", otp)

	require.NoError(t, err)

	verifiedUser, _ := repo.GetUserByEmail("john@example.com")
	assert.Empty(t, verifiedUser.PasswordResetOTP)
}

func TestEmailVerificationService_VerifyPasswordResetOTP_WrongOTP(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 5, 3)

	user := &authdomain.User{
		ID:    uuid.New().String(),
		Name:  "John",
		Email: "john@example.com",
	}
	repo.CreateUser(user)

	service.SendPasswordResetOTP("john@example.com")

	err := service.VerifyPasswordResetOTP("john@example.com", "000000")

	assert.Equal(t, ErrInvalidResetToken, err)
}

func TestEmailVerificationService_VerifyPasswordResetOTP_MaxAttempts(t *testing.T) {
	repo := NewMockUserRepositoryForOTP()
	service := NewEmailVerificationService(repo, 5, 3)

	user := &authdomain.User{
		ID:    uuid.New().String(),
		Name:  "John",
		Email: "john@example.com",
	}
	repo.CreateUser(user)

	service.SendPasswordResetOTP("john@example.com")

	err1 := service.VerifyPasswordResetOTP("john@example.com", "000000")
	err2 := service.VerifyPasswordResetOTP("john@example.com", "000000")
	err3 := service.VerifyPasswordResetOTP("john@example.com", "000000")

	assert.Equal(t, ErrInvalidResetToken, err1)
	assert.Equal(t, ErrInvalidResetToken, err2)
	assert.Equal(t, ErrMaxResetAttempts, err3)
}
