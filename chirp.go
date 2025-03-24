package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}
	type responseBody struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := requestBody{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding the incoming Chirp")
		return
	}
	const maxChirpLength = 140
	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	respondWithJSON(w, http.StatusOK, responseBody{Valid: true})
}
