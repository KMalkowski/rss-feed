package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/KMalkowski/rss-feed/internal/database"
	"github.com/google/uuid"
)

type AuthedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (a *apiConfig) authMiddleware(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(auth) < 2 {
			respondWithError(w, http.StatusForbidden, "authenticate yourself to get the user data")
			return
		}

		user, err := a.DB.GetUserByApiKey(r.Context(), strings.Replace(auth, "Bearer ", "", 1))
		if err != nil {
			log.Fatalln(err.Error(), auth)
			respondWithError(w, http.StatusInternalServerError, "something went wront try again later")
			return
		}

		if user == (database.User{}) {
			respondWithError(w, http.StatusNotFound, "could not found the user in the database")
			return
		}

		handler(w, r, user)
	}
}

func (a *apiConfig) HandleCreateUser(w http.ResponseWriter, req *http.Request) {
	type createUserBody struct {
		Name string `json:"name"`
	}

	body := createUserBody{}
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		respondWithError(w, 400, "wrong request body")
		return
	}

	user, err := a.DB.CreateUser(req.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      body.Name,
	})

	if err != nil {
		log.Fatalln(err.Error())
		respondWithError(w, 500, "count not create user, try again later")
		return
	}

	respondWithJson(w, 201, user)
}

func (a *apiConfig) GetUserByApiKeyHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJson(w, 200, user)
}
