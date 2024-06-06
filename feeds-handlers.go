package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/KMalkowski/rss-feed/internal/database"
	"github.com/google/uuid"
)

func (a *apiConfig) CreateFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type CreateFeedBody struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	body := CreateFeedBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "provide a name and url to create a new feed")
		return
	}

	feed, err := a.DB.CreateFeed(r.Context(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      body.Name,
			Url:       body.Url,
			UserID:    user.ID,
		})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not save new feed to the database")
		return
	}

	follow, err := a.DB.CreateFollow(r.Context(), database.CreateFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})

	if err != nil {
		log.Fatalln(err.Error())
		respondWithError(w, http.StatusInternalServerError, "could not save new feed follow to the database")
		return
	}

	type CreateFeedHandlerResponse struct {
		Feed   database.Feed
		Follow database.FeedFollow
	}
	respondWithJson(w, http.StatusCreated, CreateFeedHandlerResponse{Feed: feed, Follow: follow})
}

func (a *apiConfig) GetFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := a.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not query any feeds, sry")
		return
	}

	respondWithJson(w, http.StatusOK, feeds)
}
