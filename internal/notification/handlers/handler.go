package notificationhandlers

import (
	"net/http"
	"strconv"

	notificationservices "github.com/suryansh74/chat_app/internal/notification/services"
	"github.com/suryansh74/chat_app/shared/helper"
	"github.com/suryansh74/chat_app/shared/middleware"
	"github.com/suryansh74/chat_app/shared/token"
)

type NotificationHandler struct {
	notificationService notificationservices.NotificationServicePort
}

func NewNotificationHandler(notificationService notificationservices.NotificationServicePort) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok || payload == nil || payload.User == nil {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	userID := payload.User.ID

	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil {
			limit = n
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if n, err := strconv.Atoi(o); err == nil {
			offset = n
		}
	}

	notifications, err := h.notificationService.GetNotifications(userID, limit, offset)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{"notifications": notifications})
}

func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok || payload == nil || payload.User == nil {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	userID := payload.User.ID

	count, err := h.notificationService.GetUnreadCount(userID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{"unread_count": count})
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	notificationID := r.URL.Query().Get("id")
	if notificationID == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
		return
	}

	err := h.notificationService.MarkAsRead(notificationID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok || payload == nil || payload.User == nil {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	userID := payload.User.ID

	err := h.notificationService.MarkAllAsRead(userID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "all marked as read"})
}
