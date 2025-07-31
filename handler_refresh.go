package main

import (
	"net/http"
	"time"

	"github.com/GircysRomualdas/chirpycli/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	tokenJWT, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get JWT token", err)
		return
	}
	refresh_token, err := cfg.db.GetRefreshToken(r.Context(), tokenJWT)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}
	if refresh_token.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", nil)
		return
	}
	if refresh_token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked", nil)
		return
	}
	userID := refresh_token.UserID
	newJWT, err := auth.MakeJWT(userID, cfg.JWTSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: newJWT,
	})
}
