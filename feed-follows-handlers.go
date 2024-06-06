package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/KMalkowski/rss-feed/internal/database"
	"github.com/google/uuid"
)

func (a *apiConfig) CreateFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type CreateFeedFollowBody struct {
		FeedId string `json:"feed_id"`
	}

	body := CreateFeedFollowBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "provide a feed id to follow one")
		return
	}

	followUuid, err := uuid.FromBytes([]byte(body.FeedId))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get a feed uuid")
		return
	}

	follow, err := a.DB.CreateFollow(r.Context(), database.CreateFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    followUuid,
	})

	respondWithJson(w, http.StatusCreated, follow)
}
