package repositories

import (
	"errors"
	"sync"

	apperr "github.com/suryansh74/chat_app/internal/auth/apperr"
	authdomain "github.com/suryansh74/chat_app/internal/auth/domain"
)

type InMemoryUserRepository struct {
	mu           sync.RWMutex
	usersByID    map[string]*authdomain.User
	usersByEmail map[string]*authdomain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	repo := &InMemoryUserRepository{
		usersByID:    make(map[string]*authdomain.User),
		usersByEmail: make(map[string]*authdomain.User),
	}
	repo.seed()
	return repo
}

func (r *InMemoryUserRepository) EmailExists(email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.usersByEmail[email]
	return exists, nil
}

func (r *InMemoryUserRepository) CreateUser(user *authdomain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.usersByEmail[user.Email]; exists {
		return apperr.NewEmailAlreadyExists(user.Email)
	}

	r.usersByID[user.ID] = user
	r.usersByEmail[user.Email] = user

	return nil
}

func (r *InMemoryUserRepository) GetUserByEmail(email string) (*authdomain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.usersByEmail[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *InMemoryUserRepository) seed() {
	seedUsers := []*authdomain.User{
		{ID: "00000000-0000-0000-0000-000000000001", Name: "Alice", Email: "alice@example.com", Password: "Password1!"},
		{ID: "00000000-0000-0000-0000-000000000002", Name: "Bob", Email: "bob@example.com", Password: "Password2!"},
	}

	for _, user := range seedUsers {
		r.usersByID[user.ID] = user
		r.usersByEmail[user.Email] = user
	}
}
