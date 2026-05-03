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

type MockEmailVerificationService struct {
	mock.Mock
}

func (m *MockEmailVerificationService) SendOTP(email string) (string, error) {
	args := m.Called(email)
	return args.String(0), args.Error(1)
}

func (m *MockEmailVerificationService) VerifyOTP(email, otp string) error {
	args := m.Called(email, otp)
	return args.Error(0)
}

func (m *MockEmailVerificationService) IsVerified(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmailVerificationService) SendPasswordResetOTP(email string) (string, error) {
	args := m.Called(email)
	return args.String(0), args.Error(1)
}

func (m *MockEmailVerificationService) VerifyPasswordResetOTP(email, otp string) error {
	args := m.Called(email, otp)
	return args.Error(0)
}

func (m *MockEmailVerificationService) SetPassword(email, newPassword string) error {
	args := m.Called(email, newPassword)
	return args.Error(0)
}

type MockEmailSender struct {
	mock.Mock
	sentEmails []EmailRecord
}

type EmailRecord struct {
	To      string
	Subject string
	Body    string
}

func (m *MockEmailSender) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	m.sentEmails = append(m.sentEmails, EmailRecord{
		To:      to,
		Subject: subject,
		Body:    body,
	})
	return args.Error(0)
}

func createUnverifiedUserContext() context.Context {
	user := &token.Payload{
		User: &token.TokenUser{
			ID:         "user-id",
			Name:       "John",
			Email:      "john@example.com",
			IsVerified: false,
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Hour),
	}
	return context.WithValue(context.Background(), middleware.UserContextKey, user)
}

func createVerifiedUserContext() context.Context {
	user := &token.Payload{
		User: &token.TokenUser{
			ID:         "user-id",
			Name:       "John",
			Email:      "john@example.com",
			IsVerified: true,
		},
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Hour),
	}
	return context.WithValue(context.Background(), middleware.UserContextKey, user)
}

func TestEmailVerificationHandler_SendOTP_Success(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	mockService.On("SendOTP", "john@example.com").Return("123456", nil)
	mockEmailSender.On("SendEmail", "john@example.com", mock.Anything, mock.Anything).Return(nil)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/send_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createUnverifiedUserContext())
	w := httptest.NewRecorder()

	handler.SendOTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEmailVerificationHandler_SendOTP_AlreadyVerified(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/send_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createVerifiedUserContext())
	w := httptest.NewRecorder()

	handler.SendOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEmailVerificationHandler_SendOTP_Unauthorized(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/send_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.SendOTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestEmailVerificationHandler_SendOTP_ServiceError(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	mockService.On("SendOTP", "john@example.com").Return("", authservices.ErrAlreadyVerified)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/send_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createUnverifiedUserContext())
	w := httptest.NewRecorder()

	handler.SendOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEmailVerificationHandler_SendOTP_EmailSendFails(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	mockService.On("SendOTP", "john@example.com").Return("123456", nil)
	mockEmailSender.On("SendEmail", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/send_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createUnverifiedUserContext())
	w := httptest.NewRecorder()

	handler.SendOTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestEmailVerificationHandler_VerifyOTP_Success(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	mockService.On("VerifyOTP", "john@example.com", "123456").Return(nil)

	newTokenUser := &token.TokenUser{
		ID:         "user-id",
		Name:       "John",
		Email:      "john@example.com",
		IsVerified: true,
	}
	mockToken.On("CreateToken", newTokenUser, mock.AnythingOfType("time.Duration")).Return("new-valid-token", nil)

	body := `{"otp":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createUnverifiedUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEmailVerificationHandler_VerifyOTP_InvalidOTP(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	mockService.On("VerifyOTP", "john@example.com", "123456").Return(authservices.ErrInvalidOTP)

	body := `{"otp":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createUnverifiedUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEmailVerificationHandler_VerifyOTP_ExpiredOTP(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	mockService.On("VerifyOTP", "john@example.com", "123456").Return(authservices.ErrOTPExpired)

	body := `{"otp":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createUnverifiedUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEmailVerificationHandler_VerifyOTP_MaxAttempts(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	mockService.On("VerifyOTP", "john@example.com", "123456").Return(authservices.ErrMaxAttempts)

	body := `{"otp":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createUnverifiedUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEmailVerificationHandler_VerifyOTP_EmptyOTP(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	body := `{"otp":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createUnverifiedUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEmailVerificationHandler_VerifyOTP_InvalidLength(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	body := `{"otp":"12345"}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createUnverifiedUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEmailVerificationHandler_VerifyOTP_NonNumeric(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	body := `{"otp":"12a45b"}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createUnverifiedUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEmailVerificationHandler_VerifyOTP_AlreadyVerified(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	body := `{"otp":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/verify_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createVerifiedUserContext())
	w := httptest.NewRecorder()

	handler.VerifyOTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEmailVerificationHandler_Verified_Success(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	req := httptest.NewRequest(http.MethodGet, "/api/email_verification/verified", nil)
	req = req.WithContext(createVerifiedUserContext())
	w := httptest.NewRecorder()

	handler.Verified(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, true, data["verified"])
}

func TestEmailVerificationHandler_Verified_Unverified(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	req := httptest.NewRequest(http.MethodGet, "/api/email_verification/verified", nil)
	req = req.WithContext(createUnverifiedUserContext())
	w := httptest.NewRecorder()

	handler.Verified(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, false, data["verified"])
}

func TestEmailVerificationHandler_SendOTP_OTPNotInResponse(t *testing.T) {
	mockService := new(MockEmailVerificationService)
	mockEmailSender := new(MockEmailSender)
	mockToken := new(MockTokenMaker)

	handler := NewEmailVerificationHandler(mockService, mockEmailSender, mockToken, 3600)

	mockService.On("SendOTP", "john@example.com").Return("123456", nil)
	mockEmailSender.On("SendEmail", "john@example.com", mock.Anything, mock.Anything).Return(nil)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/email_verification/send_otp", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(createUnverifiedUserContext())
	w := httptest.NewRecorder()

	handler.SendOTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotContains(t, w.Body.String(), "123456", "OTP should not be in response")
}
