package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RivellionCS/chirpy/internal/auth"
	"github.com/RivellionCS/chirpy/internal/database"
)

func(cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error getting token from header: %s", err)
		respondWithError(w, http.StatusUnauthorized, "error getting token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtKey)
	if err != nil {
		log.Printf("error validating key: %s", err)
		respondWithError(w, http.StatusUnauthorized, "error validating token")
		return
	}

	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding body: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error decoding body")
		return
	}

	newHashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error updating password")
		return
	}

	queryParams:= database.UpdateUserEmailAndPasswordParams{
		ID: userID,
		Email: params.Email,
		HashedPassword: newHashedPassword,
	}

	user, err := cfg.databaseQueries.UpdateUserEmailAndPassword(r.Context(), queryParams)
	if err != nil {
		log.Printf("error updating email and password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error updating user")
		return
	}
	
	userJSON := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w, http.StatusOK, userJSON)
}