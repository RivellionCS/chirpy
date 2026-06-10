package main

import (
	"log"
	"net/http"

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
	chirps, err := cfg.databaseQueries.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error getting all chirps: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error getting all chirps")
		return
	}
	chirpsSlice := []Chirp{}
	for _, chirp := range chirps {
		chirpJSON := Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		}
		chirpsSlice = append(chirpsSlice, chirpJSON)
	}
	respondWithJSON(w, http.StatusOK, chirpsSlice)
}