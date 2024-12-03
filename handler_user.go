package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kensonjohnson/go-rss-aggregator/internal/auth"
	"github.com/kensonjohnson/go-rss-aggregator/internal/database"
)

func (api *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type bodyParams struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := bodyParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 422, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := api.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Could not create user: %v", err))
		return
	}

	respondWithJSON(w, 201, dbUserToUser(user))
}

func (api *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Auth error: %v", err))
		return
	}

	user, err := api.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get user: %v", err))
		return
	}

	respondWithJSON(w, 200, dbUserToUser(user))
}