package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Edudlufetips1/chirpy/internal/auth"
	"github.com/Edudlufetips1/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := cfg.DBQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	valid, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error checking password")
		return
	}
	if !valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	expiresIn := time.Hour

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating refresh token")
		return
	}

	_, err = cfg.DBQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		log.Printf("Error storing refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error storing refresh token")
		return
	}
	respondWithJSON(w, http.StatusOK, struct {
		User         User   `json:"user"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}
