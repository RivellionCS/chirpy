package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RivellionCS/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding json: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error logging in")
		return
	}
	user, err := cfg.databaseQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error getting user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error getting user")
		return
	}
	check, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error checking password and hash: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error checking password")
		return
	}
	if check {
		userJSON := User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		}
		respondWithJSON(w, http.StatusOK, userJSON)
	} else {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
	}
}