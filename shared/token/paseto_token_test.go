package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testKey = "12345678901234567890123456789012"

func TestNewPasetoMaker_Success(t *testing.T) {
	maker, err := NewPasetoMaker(testKey)

	require.NoError(t, err)
	assert.NotNil(t, maker)
}

func TestNewPasetoMaker_InvalidKeySize(t *testing.T) {
	_, err := NewPasetoMaker("short-key")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid key size")
}

func TestNewPasetoMaker_EmptyKey(t *testing.T) {
	_, err := NewPasetoMaker("")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid key size")
}

func TestNewPasetoMaker_KeyTooLong(t *testing.T) {
	_, err := NewPasetoMaker("123456789012345678901234567890123")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid key size")
}

func TestPasetoMaker_CreateToken_Success(t *testing.T) {
	maker, err := NewPasetoMaker(testKey)
	require.NoError(t, err)

	user := &TokenUser{
		ID:    "test-user-id",
		Name:  "Test User",
		Email: "test@example.com",
		Image: "https://example.com/img.png",
	}

	tokenStr, err := maker.CreateToken(user, 24*time.Hour)

	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)
}

func TestPasetoMaker_CreateToken_VerifyToken(t *testing.T) {
	maker, err := NewPasetoMaker(testKey)
	require.NoError(t, err)

	user := &TokenUser{
		ID:    "test-user-id",
		Name:  "Test User",
		Email: "test@example.com",
		Image: "https://example.com/img.png",
	}

	tokenStr, err := maker.CreateToken(user, 24*time.Hour)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(tokenStr)

	require.NoError(t, err)
	assert.Equal(t, user.ID, payload.User.ID)
	assert.Equal(t, user.Name, payload.User.Name)
	assert.Equal(t, user.Email, payload.User.Email)
	assert.Equal(t, user.Image, payload.User.Image)
}

func TestPasetoMaker_VerifyToken_InvalidToken(t *testing.T) {
	maker, err := NewPasetoMaker(testKey)
	require.NoError(t, err)

	_, err = maker.VerifyToken("invalid-token-string")

	assert.Error(t, err)
}

func TestPasetoMaker_VerifyToken_EmptyToken(t *testing.T) {
	maker, err := NewPasetoMaker(testKey)
	require.NoError(t, err)

	_, err = maker.VerifyToken("")

	assert.Error(t, err)
}

func TestPasetoMaker_CreateToken_DifferentKeys(t *testing.T) {
	maker1, err := NewPasetoMaker("12345678901234567890123456789012")
	require.NoError(t, err)

	maker2, err := NewPasetoMaker("abcdefghijklmnopqrstuvwxyz123456")
	require.NoError(t, err)

	user := &TokenUser{
		ID:    "test-user-id",
		Name:  "Test User",
		Email: "test@example.com",
		Image: "https://example.com/img.png",
	}

	tokenStr, err := maker1.CreateToken(user, 24*time.Hour)
	require.NoError(t, err)

	_, err = maker2.VerifyToken(tokenStr)

	assert.Error(t, err)
}

func TestNewPayload_Success(t *testing.T) {
	user := &TokenUser{
		ID:    "test-user-id",
		Name:  "Test User",
		Email: "test@example.com",
		Image: "https://example.com/img.png",
	}

	payload, err := NewPayload(user, 24*time.Hour)

	require.NoError(t, err)
	assert.Equal(t, user, payload.User)
	assert.True(t, payload.IssuedAt.Before(time.Now().Add(time.Second)))
	assert.True(t, payload.ExpiredAt.After(time.Now().Add(23*time.Hour)))
}

func TestPayload_Valid_ExpiredToken(t *testing.T) {
	user := &TokenUser{
		ID:    "test-user-id",
		Name:  "Test User",
		Email: "test@example.com",
		Image: "https://example.com/img.png",
	}

	payload := &Payload{
		User:      user,
		IssuedAt:  time.Now().Add(-2 * time.Hour),
		ExpiredAt: time.Now().Add(-1 * time.Hour),
	}

	err := payload.Valid()

	assert.Equal(t, ErrExpiredToken, err)
}

func TestPayload_Valid_ValidToken(t *testing.T) {
	user := &TokenUser{
		ID:    "test-user-id",
		Name:  "Test User",
		Email: "test@example.com",
		Image: "https://example.com/img.png",
	}

	payload := &Payload{
		User:      user,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	err := payload.Valid()

	assert.NoError(t, err)
}

func TestPayload_Valid_JustExpired(t *testing.T) {
	user := &TokenUser{
		ID:    "test-user-id",
		Name:  "Test User",
		Email: "test@example.com",
		Image: "https://example.com/img.png",
	}

	payload := &Payload{
		User:      user,
		IssuedAt:  time.Now().Add(-2 * time.Hour),
		ExpiredAt: time.Now().Add(-1 * time.Second),
	}

	err := payload.Valid()

	assert.Equal(t, ErrExpiredToken, err)
}

func TestPayload_Valid_FutureExpiration(t *testing.T) {
	user := &TokenUser{
		ID:    "test-user-id",
		Name:  "Test User",
		Email: "test@example.com",
		Image: "https://example.com/img.png",
	}

	payload := &Payload{
		User:      user,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(365 * 24 * time.Hour),
	}

	err := payload.Valid()

	assert.NoError(t, err)
}

func TestPasetoMaker_CreateToken_ShortDuration(t *testing.T) {
	maker, err := NewPasetoMaker(testKey)
	require.NoError(t, err)

	user := &TokenUser{
		ID:    "test-user-id",
		Name:  "Test User",
		Email: "test@example.com",
		Image: "https://example.com/img.png",
	}

	tokenStr, err := maker.CreateToken(user, 1*time.Second)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(tokenStr)
	require.NoError(t, err)
	assert.Equal(t, user.ID, payload.User.ID)
}

func TestPasetoMaker_CreateToken_NilUser(t *testing.T) {
	maker, err := NewPasetoMaker(testKey)
	require.NoError(t, err)

	tokenStr, err := maker.CreateToken(nil, 24*time.Hour)

	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)
}

func TestPasetoMaker_CreateToken_EmptyFields(t *testing.T) {
	maker, err := NewPasetoMaker(testKey)
	require.NoError(t, err)

	user := &TokenUser{
		ID:    "",
		Name:  "",
		Email: "",
		Image: "",
	}

	tokenStr, err := maker.CreateToken(user, 24*time.Hour)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(tokenStr)
	require.NoError(t, err)
	assert.Empty(t, payload.User.ID)
	assert.Empty(t, payload.User.Name)
	assert.Empty(t, payload.User.Email)
	assert.Empty(t, payload.User.Image)
}
