package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"

	authservices "github.com/suryansh74/chat_app/internal/auth/services"
	"github.com/suryansh74/chat_app/shared/helper"
	"github.com/suryansh74/chat_app/shared/token"
)

type PasswordResetHandler struct {
	service     authservices.EmailVerificationServicePort
	emailSender EmailSender
	tokenMaker  token.Maker
	redirectURL string
}

func NewPasswordResetHandler(service authservices.EmailVerificationServicePort, emailSender EmailSender, tokenMaker token.Maker, redirectURL string) *PasswordResetHandler {
	return &PasswordResetHandler{
		service:     service,
		emailSender: emailSender,
		tokenMaker:  tokenMaker,
		redirectURL: redirectURL,
	}
}

func (h *PasswordResetHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if input.Email == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "email is required",
		})
		return
	}

	otp, err := h.service.SendPasswordResetOTP(input.Email)
	if err != nil {
		if err.Error() == "record not found" {
			helper.WriteJSON(w, http.StatusNotFound, map[string]string{
				"error": "email not found",
			})
			return
		}
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to send OTP",
		})
		return
	}

	subject := "Your Password Reset OTP"
	body := "Your password reset code is: " + otp + ". This code will expire in 5 minutes."

	if err := h.emailSender.SendEmail(input.Email, subject, body); err != nil {
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
	var input struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if input.Email == "" || input.OTP == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "email and OTP are required",
		})
		return
	}

	otp := input.OTP
	if len(otp) != 6 || !isNumeric(otp) {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "otp must be exactly 6 digits",
		})
		return
	}

	if err := h.service.VerifyPasswordResetOTP(input.Email, otp); err != nil {
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

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "OTP verified successfully",
	})
}

func (h *PasswordResetHandler) SetPassword(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email                string `json:"email"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"password_confirmation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if input.Email == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "email is required",
		})
		return
	}

	password := input.Password
	passwordConfirmation := input.PasswordConfirmation

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

	if err := h.service.SetPassword(input.Email, password); err != nil {
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
