package chathandlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	chatdomain "github.com/suryansh74/chat_app/internal/chat/domain"
	chatservices "github.com/suryansh74/chat_app/internal/chat/services"
	"github.com/suryansh74/chat_app/shared/helper"
	"github.com/suryansh74/chat_app/shared/middleware"
	"github.com/suryansh74/chat_app/shared/token"
)

type ChatHandler struct {
	chatService chatservices.ChatServicePort
}

func NewChatHandler(chatService chatservices.ChatServicePort) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var input chatdomain.SendMessageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok || payload == nil || payload.User == nil {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	userID := payload.User.ID

	msg, err := h.chatService.SendMessage(userID, input.ToUserID, input.Content)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{"message": msg})
}

func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	friendID := r.URL.Query().Get("friend_id")
	if friendID == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "friend_id required"})
		return
	}

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

	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok || payload == nil || payload.User == nil {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	userID := payload.User.ID

	messages, err := h.chatService.GetMessages(userID, friendID, limit, offset)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{"messages": messages})
}

func (h *ChatHandler) SearchMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "query parameter required"})
		return
	}

	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok || payload == nil || payload.User == nil {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	userID := payload.User.ID

	messages, err := h.chatService.SearchMessages(userID, query)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{"messages": messages})
}

func (h *ChatHandler) SearchConversationMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	friendID := r.URL.Query().Get("friend_id")
	if query == "" || friendID == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "query and friend_id required"})
		return
	}

	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok || payload == nil || payload.User == nil {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	userID := payload.User.ID

	messages, err := h.chatService.SearchConversationMessages(userID, friendID, query)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{"messages": messages})
}
