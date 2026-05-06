package presencehandlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/suryansh74/chat_app/shared/helper"
)

type PresenceHandler struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewPresenceHandler(redisClient *redis.Client) *PresenceHandler {
	return &PresenceHandler{redisClient: redisClient, ctx: context.Background()}
}

func (h *PresenceHandler) CheckPresence(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("user_ids")
	if ids == "" {
		helper.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "user_ids required"})
		return
	}

	onlineMap := make(map[string]bool)
	for _, id := range strings.Split(ids, ",") {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		count, _ := h.redisClient.Exists(h.ctx, "online:"+id).Result()
		onlineMap[id] = count > 0
	}

	helper.WriteJSON(w, http.StatusOK, map[string]interface{}{"online": onlineMap})
}
