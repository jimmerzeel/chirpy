package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}
	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
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
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	inputMsg := strings.Split(chirp.Body, " ")
	for i, word := range inputMsg {
		if slices.Contains(profanity, strings.ToLower(word)) {
			inputMsg[i] = "****"
		}
	}
	outputMsg := strings.Join(inputMsg, " ")
	respondWithJSON(w, http.StatusOK, responseBody{CleanedBody: outputMsg})
}
