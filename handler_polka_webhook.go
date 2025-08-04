package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GircysRomualdas/chirpycli/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID", err)
		return
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}
	if apiKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}
	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
