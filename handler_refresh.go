package main

import (
	"net/http"
	"time"

	"github.com/Edudlufetips1/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	user, err := cfg.DBQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	newToken, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": newToken,
	})
}
