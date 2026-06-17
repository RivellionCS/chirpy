package main

import (
	"log"
	"net/http"

	"github.com/RivellionCS/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		log.Printf("Error parsing chirp id: %s", err)
		respondWithError(w, http.StatusBadRequest, "error getting chirp")
		return
	}

	chirp, err := cfg.databaseQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error getting chirp: %s", err)
		respondWithError(w, http.StatusNotFound, "error getting chirp")
		return
	}

	chirpJSON := Chirp{
		ID: chirpID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirpJSON)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")

	chirpsSlice := []database.Chirp{}

	if authorID == "" {
		chirps, err := cfg.databaseQueries.GetAllChirps(r.Context())
		if err != nil {
			log.Printf("Error getting all chirps: %s", err)
			respondWithError(w, http.StatusInternalServerError, "Error getting all chirps")
			return
		} 
		chirpsSlice = chirps
	} else {
		userID, err := uuid.Parse(authorID)
		if err != nil {
			log.Printf("Error parsing userID: %s", err)
			respondWithError(w, http.StatusBadRequest, "Error parsing userID")
			return
		}

		chirps, err := cfg.databaseQueries.GetAllChirpsForUserID(r.Context(), userID)
		if err != nil {
			log.Printf("Error getting all chirps: %s", err)
			respondWithError(w, http.StatusInternalServerError, "Error getting all chirps for user")
			return
		}
		chirpsSlice = chirps
	}

	chirpsSliceReturn := []Chirp{}
	for _, chirp := range chirpsSlice {
		chirpJSON := Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		}
		chirpsSliceReturn = append(chirpsSliceReturn, chirpJSON)
	}

	respondWithJSON(w, http.StatusOK, chirpsSliceReturn)
}