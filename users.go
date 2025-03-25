package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	email := requestBody{}
	err := decoder.Decode(&email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request")
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), email.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Input is not an email address")
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}
