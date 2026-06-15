package main

import (
	"log"
	"net/http"

	"github.com/RivellionCS/chirpy/internal/auth"
)

func(cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error could not get header: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error getting token")
		return
	}

	err = cfg.databaseQueries.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		log.Printf("error revoking token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error revoking token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}