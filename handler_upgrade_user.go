package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RivellionCS/chirpy/internal/auth"
	"github.com/google/uuid"
)

func(cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	polkaApiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("error getting apikey: %s", err)
		respondWithError(w, http.StatusUnauthorized, "error getting apikey")
		return
	}

	if polkaApiKey != cfg.polkaKey {
		log.Printf("error api keys do not match: %s", err)
		respondWithError(w, http.StatusUnauthorized, "api keys do not match")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding body: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error decoding body")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Printf("error parsing userID: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error parsing userID")
		return
	}

	err = cfg.databaseQueries.UpgradeToChirpyRedByID(r.Context(), userID)
	if err != nil {
		log.Printf("error upgrading user to chirpy red: %s", err)
		respondWithError(w, http.StatusNotFound, "error could not find user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}