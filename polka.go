package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jimmerzeel/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	polkaKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Cannot find Polka API key")
		return
	}
	if cfg.polkaKey != polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Polka API key does not match")
		return
	}

	decoder := json.NewDecoder(r.Body)
	input := requestBody{}
	err = decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request")
		return
	}
	if input.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_, err = cfg.db.UpgradeUser(r.Context(), input.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Cannot find user")
			return
		}
		respondWithError(w, http.StatusNotFound, "Cannot update user")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
