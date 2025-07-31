package main

import (
	"net/http"

	"github.com/GircysRomualdas/chirpycli/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	tokenJWT, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get refresh token", err)
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), tokenJWT)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't revoke refresh token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
