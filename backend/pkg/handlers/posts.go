package handlers

import (
	"encoding/json"
	"net/http"
	"social/models"
	"social/pkg/db/sqlite"
	"social/queries"
	"social/queries/middleware"
)

// ReactToPost handles adding or removing an emoji reaction to a post
func ReactToPost(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if postID == "" {
		http.Error(w, `{"error":"missing post id"}`, http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req models.ReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Emoji == "" {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	if err := queries.ToggleReaction(sqlite.DB, postID, userID, req.Emoji); err != nil {
		http.Error(w, `{"error":"failed to process reaction"}`, http.StatusInternalServerError)
		return
	}

	// Return the updated reaction count for this post
	reactions, err := queries.GetReactionsByPostIDStr(sqlite.DB, postID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch reactions"}`, http.StatusInternalServerError)
		return
	}
	if reactions == nil {
		reactions = []models.Reaction{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"reactions": reactions,
	})
}
