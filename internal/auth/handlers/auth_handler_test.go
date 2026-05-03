package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	authdomain "github.com/suryansh74/chat_app/internal/auth/domain"
	"github.com/suryansh74/chat_app/shared/token"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(input *authdomain.RegisterInput) (*authdomain.User, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authdomain.User), args.Error(1)
}

func (m *MockAuthService) ValidateRegisterInput(input *authdomain.RegisterInput) []error {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]error)
}

func (m *MockAuthService) Login(input *authdomain.LoginInput) (*authdomain.User, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authdomain.User), args.Error(1)
}

func (m *MockAuthService) ValidateLoginInput(input *authdomain.LoginInput) []error {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]error)
}

func (m *MockAuthService) Logout() error {
	args := m.Called()
	return args.Error(0)
}

type testPayload struct {
	User      *token.TokenUser `json:"user"`
	IssuedAt  time.Time        `json:"issued_at"`
	ExpiredAt time.Time        `json:"expired_at"`
}

type MockTokenMaker struct {
	mock.Mock
}

func (m *MockTokenMaker) CreateToken(user *token.TokenUser, duration time.Duration) (string, error) {
	args := m.Called(user, duration)
	return args.String(0), args.Error(1)
}

func (m *MockTokenMaker) VerifyToken(tokenStr string) (*token.Payload, error) {
	args := m.Called(tokenStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	switch v := args.Get(0).(type) {
	case *token.Payload:
		return v, args.Error(1)
	case testPayload:
		return &token.Payload{
			User:      v.User,
			IssuedAt:  v.IssuedAt,
			ExpiredAt: v.ExpiredAt,
		}, args.Error(1)
	default:
		return nil, args.Error(1)
	}
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockService := new(MockAuthService)
	mockToken := new(MockTokenMaker)

	handler := NewAuthHandler(mockService, mockToken, 3600, "Lax")

	mockService.On("ValidateRegisterInput", mock.AnythingOfType("*authdomain.RegisterInput")).Return([]error{})

	user := &authdomain.User{
		ID:    "test-id",
		Name:  "John",
		Email: "john@example.com",
	}
	mockService.On("Register", mock.AnythingOfType("*authdomain.RegisterInput")).Return(user, nil)

	tokenUser := &token.TokenUser{
		ID:    "test-id",
		Name:  "John",
		Email: "john@example.com",
	}
	mockToken.On("CreateToken", tokenUser, mock.AnythingOfType("time.Duration")).Return("valid-token", nil)

	body := `{"name":"John","email":"john@example.com","password":"Password1!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "should return 201 status")

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "user registered successfully", data["message"])

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1, "should have one cookie")
	assert.Equal(t, "session_token", cookies[0].Name, "cookie should be session_token")
}

func TestAuthHandler_Register_ValidationError(t *testing.T) {
	mockService := new(MockAuthService)
	mockToken := new(MockTokenMaker)

	handler := NewAuthHandler(mockService, mockToken, 3600, "Lax")

	mockService.On("ValidateRegisterInput", mock.AnythingOfType("*authdomain.RegisterInput")).Return([]error{assert.AnError})

	body := `{"name":"Jo","email":"invalid-email","password":"pass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "should return 400 for validation errors")
}

func TestAuthHandler_Register_EmailExists(t *testing.T) {
	mockService := new(MockAuthService)
	mockToken := new(MockTokenMaker)

	handler := NewAuthHandler(mockService, mockToken, 3600, "Lax")

	mockService.On("ValidateRegisterInput", mock.AnythingOfType("*authdomain.RegisterInput")).Return([]error{})
	mockService.On("Register", mock.AnythingOfType("*authdomain.RegisterInput")).Return(nil, assert.AnError)

	body := `{"name":"John","email":"john@example.com","password":"Password1!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockService := new(MockAuthService)
	mockToken := new(MockTokenMaker)

	handler := NewAuthHandler(mockService, mockToken, 3600, "Lax")

	mockService.On("ValidateLoginInput", mock.AnythingOfType("*authdomain.LoginInput")).Return([]error{})

	user := &authdomain.User{
		ID:    "user-id-123",
		Name:  "John",
		Email: "john@example.com",
	}
	mockService.On("Login", mock.AnythingOfType("*authdomain.LoginInput")).Return(user, nil)

	tokenUser := &token.TokenUser{
		ID:    "user-id-123",
		Name:  "John",
		Email: "john@example.com",
	}
	mockToken.On("CreateToken", tokenUser, mock.AnythingOfType("time.Duration")).Return("valid-token", nil)

	body := `{"email":"john@example.com","password":"Password1!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "should return 200 status")

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "user logged in successfully", data["message"])

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1, "should have one cookie")
	assert.Equal(t, "session_token", cookies[0].Name, "cookie should be session_token")
}

func TestAuthHandler_Login_ValidationError(t *testing.T) {
	mockService := new(MockAuthService)
	mockToken := new(MockTokenMaker)

	handler := NewAuthHandler(mockService, mockToken, 3600, "Lax")

	mockService.On("ValidateLoginInput", mock.AnythingOfType("*authdomain.LoginInput")).Return([]error{assert.AnError})

	body := `{"email":"invalid-email","password":"pass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "should return 400 for validation errors")
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := new(MockAuthService)
	mockToken := new(MockTokenMaker)

	handler := NewAuthHandler(mockService, mockToken, 3600, "Lax")

	mockService.On("ValidateLoginInput", mock.AnythingOfType("*authdomain.LoginInput")).Return([]error{})
	mockService.On("Login", mock.AnythingOfType("*authdomain.LoginInput")).Return(nil, assert.AnError)

	body := `{"email":"john@example.com","password":"WrongPassword!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "should return 401 for invalid credentials")
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	mockService := new(MockAuthService)
	mockToken := new(MockTokenMaker)

	handler := NewAuthHandler(mockService, mockToken, 3600, "Lax")

	mockService.On("Logout").Return(nil)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Logout(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "should return 200 status")

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1, "should have one cookie")
	assert.Equal(t, "session_token", cookies[0].Name, "cookie should be session_token")
	assert.Equal(t, -1, cookies[0].MaxAge, "cookie should be expired/cleared")
}
