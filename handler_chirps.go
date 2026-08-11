package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/Edudlufetips1/chirpy/internal/auth"
	"github.com/Edudlufetips1/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type ChirpRequest struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error retrieving token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		log.Printf("Error validating token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	req := ChirpRequest{}
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := filterProfanity(req.Body)

	chirp, err := cfg.DBQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorIDstring := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")

	var dbChirps []database.Chirp
	var err error

	if authorIDstring != "" {
		authorID, parseErr := uuid.Parse(authorIDstring)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
		dbChirps, err = cfg.DBQueries.GetChirpsByAuthorID(r.Context(), authorID)
	} else {
		dbChirps, err = cfg.DBQueries.GetAllChirps(r.Context())
	}

	if err != nil {
		log.Printf("Error retrieving chirps: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortParam == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(idString)
	if err != nil {
		log.Printf("Invalid chirp ID: %v", err)
		respondWithError(w, http.StatusNotFound, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DBQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error retrieving chirp: %v", err)
		respondWithError(w, http.StatusNotFound, "Could not retrieve chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error retrieving token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		log.Printf("Error validating token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	chirpIDString := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		log.Printf("Invalid chirp ID: %v", err)
		respondWithError(w, http.StatusNotFound, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DBQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error retrieving chirp: %v", err)
		respondWithError(w, http.StatusNotFound, "Could not retrieve chirp")
		return
	}

	if dbChirp.UserID != userID {
		log.Printf("User %s is not authorized to delete chirp %s", userID, chirpID)
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	err = cfg.DBQueries.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error deleting chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could not delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
