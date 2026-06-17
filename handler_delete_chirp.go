package main

import (
	"log"
	"net/http"

	"github.com/RivellionCS/chirpy/internal/auth"
	"github.com/google/uuid"
)

func(cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error getting token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "error getting token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtKey)
	if err != nil {
		log.Printf("error getting token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "error getting token")
		return
	}

	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		log.Printf("error parsing chirp ID: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error parsing chirp ID")
		return
	}

	chirp, err := cfg.databaseQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		log.Printf("error getting chirp: %s", err)
		respondWithError(w, http.StatusNotFound, "error getting chirp")
		return
	}

	if chirp.UserID != userID {
		log.Printf("error user ID for chirp does not match: %s", err)
		respondWithError(w, http.StatusForbidden, "error chirp does not belong to user")
		return
	}

	err = cfg.databaseQueries.DeleteChirpByID(r.Context(), chirp.ID)
	if err != nil {
		log.Printf("error deleting chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error could not delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}