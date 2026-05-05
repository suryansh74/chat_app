package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/suryansh74/chat_app/pkg/logger"
	"github.com/suryansh74/chat_app/shared/helper"
	"github.com/suryansh74/chat_app/shared/middleware"
	"github.com/suryansh74/chat_app/shared/token"

	authservices "github.com/suryansh74/chat_app/internal/auth/services"
)

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type EmailVerificationHandler struct {
	service      authservices.EmailVerificationServicePort
	emailSender  EmailSender
	tokenMaker   token.Maker
	cookieMaxAge int
}

func NewEmailVerificationHandler(service authservices.EmailVerificationServicePort, emailSender EmailSender, tokenMaker token.Maker, cookieMaxAge int) *EmailVerificationHandler {
	return &EmailVerificationHandler{
		service:      service,
		emailSender:  emailSender,
		tokenMaker:   tokenMaker,
		cookieMaxAge: cookieMaxAge,
	}
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (h *EmailVerificationHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	if payload.User.IsVerified {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "email already verified",
		})
		return
	}

	otp, err := h.service.SendOTP(payload.User.Email)
	if err != nil {
		logger.Log.Error("SendOTP: failed to generate OTP", "error", err.Error(), "email", payload.User.Email)
		if err == authservices.ErrAlreadyVerified {
			helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "email already verified",
			})
			return
		}
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to send OTP",
		})
		return
	}

	subject := "Your OTP for Email Verification"
	body := "Your verification code is: " + otp + ". This code will expire in 5 minutes."

	logger.Log.Info("SendOTP: sending email", "email", payload.User.Email, "otp", otp)

	if err := h.emailSender.SendEmail(payload.User.Email, subject, body); err != nil {
		logger.Log.Error("SendOTP: failed to send email", "error", err.Error(), "email", payload.User.Email)
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to send email",
		})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "OTP sent successfully",
	})
}

func (h *EmailVerificationHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	if payload.User.IsVerified {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "email already verified",
		})
		return
	}

	var input struct {
		OTP string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if input.OTP == "" {
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

	if len(input.OTP) != 6 || !isNumeric(input.OTP) {
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

	if err := h.service.VerifyOTP(payload.User.Email, input.OTP); err != nil {
		if err == authservices.ErrInvalidOTP {
			helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid OTP",
			})
			return
		}
		if err == authservices.ErrOTPExpired {
			helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "OTP has expired",
			})
			return
		}
		if err == authservices.ErrMaxAttempts {
			helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "max attempts exceeded",
			})
			return
		}
		if err == authservices.ErrAlreadyVerified {
			helper.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "email already verified",
			})
			return
		}
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to verify OTP",
		})
		return
	}

	tokenUser := &token.TokenUser{
		ID:         payload.User.ID,
		Name:       payload.User.Name,
		Email:      payload.User.Email,
		IsVerified: true,
	}

	accessToken, err := h.tokenMaker.CreateToken(tokenUser, payload.ExpiredAt.Sub(payload.IssuedAt))
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
		SameSite: 1, // SameSiteLax
	}

	http.SetCookie(w, httpCookie)

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "email verified successfully",
	})
}

func (h *EmailVerificationHandler) Verified(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("EmailVerificationHandler.Verified: called")

	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok {
		logger.Log.Error("EmailVerificationHandler.Verified: no payload in context")
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	logger.Log.Info("EmailVerificationHandler.Verified: returning verified status", "verified", payload.User.IsVerified, "email", payload.User.Email)

	helper.WriteJSON(w, http.StatusOK, map[string]bool{
		"verified": payload.User.IsVerified,
	})
}
