package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/RivellionCS/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters")
		return
	}
	params.Body = getCleanedBody(params.Body)
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	chirpParams := database.CreateChirpParams{
		Body: params.Body,
		UserID: params.UserID,
	}
	chirp, err := cfg.databaseQueries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp")
		return
	}
	chirpJSON := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, chirpJSON)
}


func getCleanedBody(body string) string {
	profanities := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	strArray := strings.Split(body, " ")
	for i, word := range strArray {
		_, ok := profanities[strings.ToLower(word)]
		if ok {
			strArray[i] = "****"
		}
	}
	cleanedBody := strings.Join(strArray, " ")
	return cleanedBody
}