package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	govalidator "github.com/go-playground/validator/v10"
	"github.com/suryansh74/chat_app/internal/auth/apperr"
	authdomain "github.com/suryansh74/chat_app/internal/auth/domain"
	authservices "github.com/suryansh74/chat_app/internal/auth/services"
	"github.com/suryansh74/chat_app/shared/helper"
	"github.com/suryansh74/chat_app/shared/middleware"
	"github.com/suryansh74/chat_app/shared/token"
	"github.com/suryansh74/chat_app/shared/validator"
)

func translateServiceError(field, message string) string {
	if strings.Contains(message, "required") {
		return fmt.Sprintf("%s is required", field)
	}
	if strings.Contains(message, "email") {
		return fmt.Sprintf("%s must be a valid email", field)
	}
	if strings.Contains(message, "min") {
		return fmt.Sprintf("%s must be at least 3 characters", field)
	}
	if strings.Contains(message, "Password") && strings.Contains(message, "uppercase") {
		return "password must contain at least one uppercase letter, one number, and one special character"
	}
	if strings.Contains(message, "match") {
		return "password and password_confirmation must match"
	}

	return message
}

type AuthHandler struct {
	service        authservices.AuthServicePort
	tokenMaker     token.Maker
	cookieMaxAge   int
	cookieSameSite string
}

func NewAuthHandler(service authservices.AuthServicePort, tokenMaker token.Maker, cookieMaxAge int, cookieSameSite string) *AuthHandler {
	return &AuthHandler{
		service:        service,
		tokenMaker:     tokenMaker,
		cookieMaxAge:   cookieMaxAge,
		cookieSameSite: cookieSameSite,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input authdomain.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	validationErrors := h.service.ValidateRegisterInput(&input)
	if len(validationErrors) > 0 {
		var formattedErrors []map[string]string

		for _, err := range validationErrors {
			if ve, ok := err.(govalidator.ValidationErrors); ok {
				translated := validator.TranslateValidationErrors(ve)
				for _, t := range translated {
					formattedErrors = append(formattedErrors, map[string]string{
						"field":   t.Field,
						"message": t.Message,
					})
				}
			} else if authErr, ok := err.(authservices.ValidationError); ok {
				formattedErrors = append(formattedErrors, map[string]string{
					"field":   authErr.Field,
					"message": translateServiceError(authErr.Field, authErr.Message),
				})
			} else {
				formattedErrors = append(formattedErrors, map[string]string{
					"field":   "error",
					"message": err.Error(),
				})
			}
		}

		helper.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": formattedErrors,
		})
		return
	}

	user, err := h.service.Register(&input)
	if err != nil {
		if apperr.IsEmailAlreadyExists(err) {
			helper.WriteJSON(w, http.StatusConflict, map[string]string{
				"error": err.Error(),
			})
			return
		}
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	tokenUser := &token.TokenUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	accessToken, err := h.tokenMaker.CreateToken(tokenUser, time.Hour)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create token",
		})
		return
	}

	httpCookie := &http.Cookie{
		Name:     "session_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   h.cookieMaxAge,
		HttpOnly: true,
		SameSite: 1,
	}

	http.SetCookie(w, httpCookie)

	helper.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "user registered successfully",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input authdomain.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	validationErrors := h.service.ValidateLoginInput(&input)
	if len(validationErrors) > 0 {
		var formattedErrors []map[string]string

		for _, err := range validationErrors {
			if ve, ok := err.(govalidator.ValidationErrors); ok {
				translated := validator.TranslateValidationErrors(ve)
				for _, t := range translated {
					formattedErrors = append(formattedErrors, map[string]string{
						"field":   t.Field,
						"message": t.Message,
					})
				}
			} else if authErr, ok := err.(authservices.ValidationError); ok {
				formattedErrors = append(formattedErrors, map[string]string{
					"field":   authErr.Field,
					"message": translateServiceError(authErr.Field, authErr.Message),
				})
			} else {
				formattedErrors = append(formattedErrors, map[string]string{
					"field":   "error",
					"message": err.Error(),
				})
			}
		}

		helper.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": formattedErrors,
		})
		return
	}

	user, err := h.service.Login(&input)
	if err != nil {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid credentials",
		})
		return
	}

	tokenUser := &token.TokenUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	duration := time.Hour
	if input.RememberMe {
		duration = 24 * time.Hour * 30
	}

	accessToken, err := h.tokenMaker.CreateToken(tokenUser, duration)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create token",
		})
		return
	}

	httpCookie := &http.Cookie{
		Name:     "session_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   h.cookieMaxAge,
		HttpOnly: true,
		SameSite: 1,
	}

	http.SetCookie(w, httpCookie)

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user logged in successfully",
	})
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]string{
			"id":    payload.User.ID,
			"name":  payload.User.Name,
			"email": payload.User.Email,
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Logout(); err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "logout failed",
		})
		return
	}

	httpCookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: 1,
	}

	http.SetCookie(w, httpCookie)

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}
