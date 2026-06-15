package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/RivellionCS/chirpy/internal/auth"
	"github.com/RivellionCS/chirpy/internal/database"
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
		respondWithError(w, http.StatusBadRequest, "error logging in")
		return
	}

	user, err := cfg.databaseQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error getting user: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	check, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error checking password and hash: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if check {
		tokenTimeLimit := time.Hour
		token, err := auth.MakeJWT(user.ID, cfg.jwtKey, time.Duration(tokenTimeLimit) * time.Second)
		if err != nil {
			log.Printf("Error creating jwt token: %s", err)
			respondWithError(w, http.StatusInternalServerError, "Error creating token")
			return
		}

		refreshTokenStr := auth.MakeRefreshToken()
		refreshTokenParams := database.CreateRefreshTokenParams{
			Token: refreshTokenStr,
			UserID: user.ID,
		}

		_, err = cfg.databaseQueries.CreateRefreshToken(r.Context(), refreshTokenParams)
		if err != nil {
			log.Printf("Error creating refresh token: %s", err)
			respondWithError(w, http.StatusInternalServerError, "Error creating token")
			return
		}

		userJSON := User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
			Token: token,
			RefreshToken: refreshTokenStr,
		}
		respondWithJSON(w, http.StatusOK, userJSON)
	} else {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
	}
}