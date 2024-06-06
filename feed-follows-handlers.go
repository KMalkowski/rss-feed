package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/KMalkowski/rss-feed/internal/database"
	"github.com/google/uuid"
)

func (a *apiConfig) CreateFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type CreateFeedFollowBody struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	body := CreateFeedFollowBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "provide a feed id to follow one")
		return
	}

	follow, err := a.DB.CreateFollow(r.Context(), database.CreateFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    body.FeedId,
	})

	respondWithJson(w, http.StatusCreated, follow)
}

func (a *apiConfig) DeleteFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type DeleteFeedFollowBody struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	body := DeleteFeedFollowBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "provide a feed id to unfollow one")
		return
	}

	rows, err := a.DB.DeleteFollow(r.Context(), database.DeleteFollowParams{FeedID: body.FeedId, UserID: user.ID})
	if err != nil {
		log.Fatalln(err.Error())
		respondWithError(w, http.StatusInternalServerError, "we could not delete the follow, try again later")
		return
	}

	if rows < 1 {
		log.Printf("%v", rows)
		respondWithError(w, http.StatusNotFound, "there is no such follow in the database")
		return
	}

	respondWithJson(w, http.StatusOK, DeleteFeedFollowBody{FeedId: body.FeedId})
}
