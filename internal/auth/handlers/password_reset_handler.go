package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	authservices "github.com/suryansh74/chat_app/internal/auth/services"
	"github.com/suryansh74/chat_app/shared/email"
	"github.com/suryansh74/chat_app/shared/helper"
	"github.com/suryansh74/chat_app/shared/middleware"
	"github.com/suryansh74/chat_app/shared/token"
)

type PasswordResetHandler struct {
	service     authservices.EmailVerificationServicePort
	emailSender *email.Sender
	tokenMaker  token.Maker
	redirectURL string
}

func NewPasswordResetHandler(service authservices.EmailVerificationServicePort, emailSender *email.Sender, tokenMaker token.Maker, redirectURL string) *PasswordResetHandler {
	return &PasswordResetHandler{
		service:     service,
		emailSender: emailSender,
		tokenMaker:  tokenMaker,
		redirectURL: redirectURL,
	}
}

func (h *PasswordResetHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	otp, err := h.service.SendPasswordResetOTP(payload.User.Email)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to send OTP",
		})
		return
	}

	subject := "Your Password Reset OTP"
	body := "Your password reset code is: " + otp + ". This code will expire in 5 minutes."

	if err := h.emailSender.SendEmail(payload.User.Email, subject, body); err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to send email",
		})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "OTP sent successfully",
	})
}

func (h *PasswordResetHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	var rawInput map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&rawInput); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	otpValue, exists := rawInput["otp"]
	if !exists {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": []map[string]string{
				{
					"field":   "otp",
					"message": "otp is required",
				},
			},
		})
		return
	}

	var otp string
	switch v := otpValue.(type) {
	case string:
		otp = v
	case float64:
		otp = fmt.Sprintf("%.0f", v)
	default:
		helper.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": []map[string]string{
				{
					"field":   "otp",
					"message": "otp must be a string or number",
				},
			},
		})
		return
	}

	if otp == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": []map[string]string{
				{
					"field":   "otp",
					"message": "otp is required",
				},
			},
		})
		return
	}

	if len(otp) != 6 || !isNumeric(otp) {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": []map[string]string{
				{
					"field":   "otp",
					"message": "otp must be exactly 6 digits",
				},
			},
		})
		return
	}

	if err := h.service.VerifyPasswordResetOTP(payload.User.Email, otp); err != nil {
		if err == authservices.ErrInvalidResetToken {
			helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid OTP",
			})
			return
		}
		if err == authservices.ErrResetTokenExpired {
			helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "OTP has expired",
			})
			return
		}
		if err == authservices.ErrMaxResetAttempts {
			helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "max attempts exceeded",
			})
			return
		}
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to verify OTP",
		})
		return
	}

	resetToken, err := h.tokenMaker.CreateToken(&token.TokenUser{
		ID:    payload.User.ID,
		Email: payload.User.Email,
	}, 5*time.Minute)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to generate reset token",
		})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"redirect":    h.redirectURL,
		"reset_token": resetToken,
	})
}

func (h *PasswordResetHandler) SetPassword(w http.ResponseWriter, r *http.Request) {
	rawInput := make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&rawInput); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	resetToken, ok := rawInput["reset_token"].(string)
	if !ok || resetToken == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "reset token is required",
		})
		return
	}

	tokenPayload, err := h.tokenMaker.VerifyToken(resetToken)
	if err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid or expired reset token",
		})
		return
	}

	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	if tokenPayload.User.ID != payload.User.ID || tokenPayload.User.Email != payload.User.Email {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid reset token",
		})
		return
	}

	password, _ := rawInput["password"].(string)
	passwordConfirmation, _ := rawInput["password_confirmation"].(string)

	if password == "" || passwordConfirmation == "" {
		var errors []map[string]string
		if password == "" {
			errors = append(errors, map[string]string{
				"field":   "password",
				"message": "password is required",
			})
		}
		if passwordConfirmation == "" {
			errors = append(errors, map[string]string{
				"field":   "password_confirmation",
				"message": "password confirmation is required",
			})
		}
		helper.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": errors,
		})
		return
	}

	if password != passwordConfirmation {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": []map[string]string{
				{
					"field":   "password_confirmation",
					"message": "password and confirmation do not match",
				},
			},
		})
		return
	}

	hasUpper, hasDigit, hasSpecial := false, false, false
	for _, ch := range password {
		switch {
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= '0' && ch <= '9':
			hasDigit = true
		case ch == '!' || ch == '@' || ch == '#' || ch == '$' || ch == '%' || ch == '^' || ch == '&' || ch == '*':
			hasSpecial = true
		}
	}
	if !(hasUpper && hasDigit && hasSpecial) {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": []map[string]string{
				{
					"field":   "password",
					"message": "password must contain at least one uppercase letter, one number, and one special character",
				},
			},
		})
		return
	}

	if err := h.service.SetPassword(payload.User.Email, password); err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to set password",
		})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "password updated successfully",
	})
}

func getReadableErrorMessage(field, tag string) string {
	field = toSentenceCase(field)
	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return "invalid email format"
	case "min":
		return field + " must be at least 8 characters"
	case "max":
		return field + " must be at most 32 characters"
	default:
		return field + " is invalid"
	}
}

func toSentenceCase(s string) string {
	if len(s) == 0 {
		return s
	}
	s = strings.ToLower(s)
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
