package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/RivellionCS/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	databaseQueries *database.Queries
	platform string
}

func main() {
	buildAndStartServer()
}

func buildAndStartServer() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return
	}
	dbQueries := database.New(db)
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		databaseQueries: dbQueries,
		platform: platform,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", myHandlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirps)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUsers)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpById)
	server := http.Server{
		Handler: mux,
		Addr: ":" + port,
	}
	log.Printf("Serving files from %s on port: %s", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

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