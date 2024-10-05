package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/adityakumaaar/gowebserver/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	type params struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	param := params{}
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error in parsing JSON: %v", err))
		return
	}

	feedfollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    param.FeedID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error in creating Feed Follow Entry %v", err))
		return
	}

	respondWithJSON(w, 201, feedfollow)
}

func (apiCfg *apiConfig) handlerGetFeedsUserFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiCfg.DB.GetFeedsUserFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error in getting feeds follows %v", err))
		return
	}
	respondWithJSON(w, 200, feeds)
}
