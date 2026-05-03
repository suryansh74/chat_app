package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	authservices "github.com/suryansh74/chat_app/internal/auth/services"
	"github.com/suryansh74/chat_app/shared/middleware"
	"github.com/suryansh74/chat_app/shared/token"
)

func createAuthUserContext() context.Context {
	user := &token.Payload{
		User: &token.TokenUser{
			ID:    "user-id",
			Name:  "John",
			Email: "john@example.com",
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Hour),
	}
	return context.WithValue(context.Background(), middleware.UserContextKey, user)
}

func TestPasswordResetHandler_SendOTP_Success(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	mockService.On("SendPasswordResetOTP", "john@example.com").Return("123456", nil)
	mockEmailSender.On("SendEmail", "john@example.com", mock.Anything, mock.Anything).Return(nil)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/send_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.SendOTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPasswordResetHandler_SendOTP_Unauthorized(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/send_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.SendOTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPasswordResetHandler_SendOTP_ServiceError(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	mockService.On("SendPasswordResetOTP", "john@example.com").Return("", assert.AnError)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/send_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.SendOTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPasswordResetHandler_SendOTP_EmailFails(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	mockService.On("SendPasswordResetOTP", "john@example.com").Return("123456", nil)
	mockEmailSender.On("SendEmail", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/send_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.SendOTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPasswordResetHandler_VerifyOTP_Success(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	mockService.On("VerifyPasswordResetOTP", "john@example.com", "123456").Return(nil)

	resetTokenUser := &token.TokenUser{
		ID:    "user-id",
		Email: "john@example.com",
	}
	mockToken.On("CreateToken", resetTokenUser, 5*time.Minute).Return("reset-token-123", nil)

	body := `{"otp":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "https://example.com/reset", data["redirect"])
	assert.NotEmpty(t, data["reset_token"])
}

func TestPasswordResetHandler_VerifyOTP_InvalidOTP(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	mockService.On("VerifyPasswordResetOTP", "john@example.com", "123456").Return(authservices.ErrInvalidResetToken)

	body := `{"otp":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPasswordResetHandler_VerifyOTP_Expired(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	mockService.On("VerifyPasswordResetOTP", "john@example.com", "123456").Return(authservices.ErrResetTokenExpired)

	body := `{"otp":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPasswordResetHandler_VerifyOTP_MaxAttempts(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	mockService.On("VerifyPasswordResetOTP", "john@example.com", "123456").Return(authservices.ErrMaxResetAttempts)

	body := `{"otp":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPasswordResetHandler_VerifyOTP_EmptyOTP(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	body := `{"otp":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPasswordResetHandler_VerifyOTP_InvalidLength(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	body := `{"otp":"12345"}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPasswordResetHandler_SetPassword_Success(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	validResetToken := &token.Payload{
		User: &token.TokenUser{
			ID:    "user-id",
			Email: "john@example.com",
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}
	mockToken.On("VerifyToken", "reset-token-123").Return(validResetToken, nil)

	mockService.On("SetPassword", "john@example.com", "NewPassword1!").Return(nil)

	body := `{"reset_token":"reset-token-123","password":"NewPassword1!","password_confirmation":"NewPassword1!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/set_password", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.SetPassword(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPasswordResetHandler_SetPassword_MissingResetToken(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	body := `{"password":"NewPassword1!","password_confirmation":"NewPassword1!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/set_password", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.SetPassword(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPasswordResetHandler_SetPassword_InvalidResetToken(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	mockToken.On("VerifyToken", "invalid-reset-token").Return(nil, assert.AnError)

	body := `{"reset_token":"invalid-reset-token","password":"NewPassword1!","password_confirmation":"NewPassword1!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/set_password", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.SetPassword(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPasswordResetHandler_SetPassword_TokenMismatch(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	differentUserToken := &token.Payload{
		User: &token.TokenUser{
			ID:    "different-user-id",
			Email: "different@example.com",
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}
	mockToken.On("VerifyToken", "reset-token-123").Return(differentUserToken, nil)

	body := `{"reset_token":"reset-token-123","password":"NewPassword1!","password_confirmation":"NewPassword1!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/set_password", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.SetPassword(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPasswordResetHandler_SetPassword_PasswordMismatch(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	validResetToken := &token.Payload{
		User: &token.TokenUser{
			ID:    "user-id",
			Email: "john@example.com",
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}
	mockToken.On("VerifyToken", "reset-token-123").Return(validResetToken, nil)

	body := `{"reset_token":"reset-token-123","password":"NewPassword1!","password_confirmation":"DifferentPass1!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/set_password", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.SetPassword(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPasswordResetHandler_SetPassword_WeakPassword(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	validResetToken := &token.Payload{
		User: &token.TokenUser{
			ID:    "user-id",
			Email: "john@example.com",
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}
	mockToken.On("VerifyToken", "reset-token-123").Return(validResetToken, nil)

	body := `{"reset_token":"reset-token-123","password":"weak","password_confirmation":"weak"}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/set_password", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.SetPassword(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPasswordResetHandler_SetPassword_MissingPassword(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	validResetToken := &token.Payload{
		User: &token.TokenUser{
			ID:    "user-id",
			Email: "john@example.com",
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}
	mockToken.On("VerifyToken", "reset-token-123").Return(validResetToken, nil)

	body := `{"reset_token":"reset-token-123","password":"","password_confirmation":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/set_password", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.SetPassword(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPasswordResetHandler_SetPassword_ServiceError(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewPasswordResetHandler(mockService, mockEmailSender, mockToken, "https://example.com/reset")

	validResetToken := &token.Payload{
		User: &token.TokenUser{
			ID:    "user-id",
			Email: "john@example.com",
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}
	mockToken.On("VerifyToken", "reset-token-123").Return(validResetToken, nil)

	mockService.On("SetPassword", "john@example.com", "NewPassword1!").Return(assert.AnError)

	body := `{"reset_token":"reset-token-123","password":"NewPassword1!","password_confirmation":"NewPassword1!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/password_reset/set_password", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createAuthUserContext())
	w := httptest.NewRecorder()

	handler.SetPassword(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
