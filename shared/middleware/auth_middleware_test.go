package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suryansh74/chat_app/shared/token"
)

type testTokenMaker struct {
	verifyFunc func(tokenStr string) (*token.Payload, error)
}

func (m *testTokenMaker) CreateToken(user *token.TokenUser, duration time.Duration) (string, error) {
	return "test-token", nil
}

func (m *testTokenMaker) VerifyToken(tokenStr string) (*token.Payload, error) {
	return m.verifyFunc(tokenStr)
}

func TestAuthMiddleware_MissingCookie(t *testing.T) {
	mockMaker := &testTokenMaker{
		verifyFunc: func(tokenStr string) (*token.Payload, error) {
			return nil, token.ErrInvalidToken
		},
	}

	mw := AuthMiddleware(mockMaker)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var response map[string]interface{}
	err := decodeJSON(resp.Body, &response)
	require.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, data["error"], "session cookie missing")
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	expectedPayload := &token.Payload{
		User: &token.TokenUser{
			ID:    "test-user-id",
			Name:  "Test User",
			Email: "test@example.com",
			Image: "https://example.com/img.png",
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	mockMaker := &testTokenMaker{
		verifyFunc: func(tokenStr string) (*token.Payload, error) {
			return expectedPayload, nil
		},
	}

	mw := AuthMiddleware(mockMaker)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, ok := r.Context().Value(UserContextKey).(*token.Payload)
		assert.True(t, ok)
		assert.Equal(t, expectedPayload, payload)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "valid-token",
	})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockMaker := &testTokenMaker{
		verifyFunc: func(tokenStr string) (*token.Payload, error) {
			return nil, token.ErrInvalidToken
		},
	}

	mw := AuthMiddleware(mockMaker)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "invalid-token",
	})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var response map[string]interface{}
	err := decodeJSON(resp.Body, &response)
	require.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, data["error"], "invalid session")
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	mockMaker := &testTokenMaker{
		verifyFunc: func(tokenStr string) (*token.Payload, error) {
			return nil, token.ErrExpiredToken
		},
	}

	mw := AuthMiddleware(mockMaker)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "expired-token",
	})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var response map[string]interface{}
	err := decodeJSON(resp.Body, &response)
	require.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, data["error"], "invalid session")
}

func TestAuthMiddleware_EmptyCookieValue(t *testing.T) {
	mockMaker := &testTokenMaker{
		verifyFunc: func(tokenStr string) (*token.Payload, error) {
			return nil, token.ErrInvalidToken
		},
	}

	mw := AuthMiddleware(mockMaker)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "",
	})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMiddleware_UserContextKey(t *testing.T) {
	assert.Equal(t, contextKey("user"), UserContextKey)
}

func decodeJSON(body io.Reader, v interface{}) error {
	return json.NewDecoder(body).Decode(v)
}
