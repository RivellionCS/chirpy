package main

import (
	"log"
	"net/http"
	"time"

	"github.com/RivellionCS/chirpy/internal/auth"
)

func(cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting token from header: %s", err)
		respondWithError(w, http.StatusUnauthorized, "error getting token")
		return
	}

	user, err := cfg.databaseQueries.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		log.Printf("Error getting user by token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "error getting token")
		return
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.jwtKey, time.Hour)
	if err != nil {
		log.Printf("Error creating token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "error creating token")
		return
	}

	params := parameters{
		Token: jwt,
	}

	respondWithJSON(w, http.StatusOK, params)
}