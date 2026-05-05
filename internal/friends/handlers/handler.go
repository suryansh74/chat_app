package friendshandlers

import (
	"encoding/json"
	"net/http"

	friendsdomain "github.com/suryansh74/chat_app/internal/friends/domain"
	friendsservices "github.com/suryansh74/chat_app/internal/friends/services"
	"github.com/suryansh74/chat_app/shared/helper"
	"github.com/suryansh74/chat_app/shared/middleware"
	"github.com/suryansh74/chat_app/shared/token"
)

type FriendsHandler struct {
	friendsService friendsservices.FriendsServicePort
}

func NewFriendsHandler(friendsService friendsservices.FriendsServicePort) *FriendsHandler {
	return &FriendsHandler{friendsService: friendsService}
}

func (h *FriendsHandler) GetFriends(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok || payload == nil || payload.User == nil {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	userID := payload.User.ID

	friends, err := h.friendsService.GetFriends(userID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{"friends": friends})
}

func (h *FriendsHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	var input friendsdomain.SendFriendRequestInput
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

	err := h.friendsService.SendFriendRequest(userID, input.ToUserID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "friend request sent"})
}

func (h *FriendsHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	var input friendsdomain.AcceptFriendRequestInput
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

	err := h.friendsService.AcceptFriendRequest(input.RequestID, userID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "friend request accepted"})
}

func (h *FriendsHandler) RejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	var input friendsdomain.AcceptFriendRequestInput
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

	err := h.friendsService.RejectFriendRequest(input.RequestID, userID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "friend request rejected"})
}

func (h *FriendsHandler) SearchByEmail(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "query parameter required"})
		return
	}

	users, err := h.friendsService.SearchByEmail(query)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{"users": users})
}

func (h *FriendsHandler) SearchFriends(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok || payload == nil || payload.User == nil {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	userID := payload.User.ID

	friends, err := h.friendsService.SearchFriends(userID, query)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{"friends": friends})
}

func (h *FriendsHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	friendID := r.URL.Query().Get("friend_id")
	if friendID == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "friend_id required"})
		return
	}

	payload, ok := r.Context().Value(middleware.UserContextKey).(*token.Payload)
	if !ok || payload == nil || payload.User == nil {
		helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	userID := payload.User.ID

	err := h.friendsService.RemoveFriend(userID, friendID)
	if err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "friend removed"})
}
