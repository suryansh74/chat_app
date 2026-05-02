package repositories

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	apperr "github.com/suryansh74/chat_app/internal/auth/apperr"
	authdomain "github.com/suryansh74/chat_app/internal/auth/domain"
)

func TestInMemoryUserRepository_EmailExists_ReturnsFalseForNonExistentEmail(t *testing.T) {
	repo := NewInMemoryUserRepository()

	exists, err := repo.EmailExists("nonexistent@example.com")

	assert.NoError(t, err)
	assert.False(t, exists, "should return false for non-existent email")
}

func TestInMemoryUserRepository_EmailExists_ReturnsTrueForExistingEmail(t *testing.T) {
	repo := NewInMemoryUserRepository()

	user := &authdomain.User{
		ID:       uuid.New().String(),
		Name:     "John",
		Email:    "john@example.com",
		Password: "Password1!",
	}
	repo.CreateUser(user)

	exists, err := repo.EmailExists("john@example.com")

	assert.NoError(t, err)
	assert.True(t, exists, "should return true for existing email")
}

func TestInMemoryUserRepository_CreateUser_AddsUserToRepository(t *testing.T) {
	repo := NewInMemoryUserRepository()
	userID := uuid.New().String()
	user := &authdomain.User{
		ID:       userID,
		Name:     "John",
		Email:    "john@example.com",
		Password: "Password1!",
	}

	err := repo.CreateUser(user)

	assert.NoError(t, err)
	exists, _ := repo.EmailExists("john@example.com")
	assert.True(t, exists, "user should be added to repository")
}

func TestInMemoryUserRepository_CreateUser_ReturnsErrorForDuplicateEmail(t *testing.T) {
	repo := NewInMemoryUserRepository()
	userID1 := uuid.New().String()
	userID2 := uuid.New().String()

	user1 := &authdomain.User{
		ID:       userID1,
		Name:     "John",
		Email:    "john@example.com",
		Password: "Password1!",
	}
	repo.CreateUser(user1)

	user2 := &authdomain.User{
		ID:       userID2,
		Name:     "Jane",
		Email:    "john@example.com",
		Password: "Password2!",
	}
	err := repo.CreateUser(user2)

	assert.Error(t, err)
	assert.True(t, apperr.IsEmailAlreadyExists(err), "should return email already exists error")
}

func TestInMemoryUserRepository_SeededUsers_ArePresent(t *testing.T) {
	repo := NewInMemoryUserRepository()

	exists, err := repo.EmailExists("alice@example.com")
	assert.NoError(t, err)
	assert.True(t, exists, "seeded user alice should exist")

	exists, err = repo.EmailExists("bob@example.com")
	assert.NoError(t, err)
	assert.True(t, exists, "seeded user bob should exist")
}
