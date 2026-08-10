package main

import (
	"net/http"

	"github.com/Edudlufetips1/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	err = cfg.DBQueries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking token")
		return
	}

	respondWithJSON(w, http.StatusNoContent, map[string]string{
		"message": "Token revoked successfully",
	})
}
