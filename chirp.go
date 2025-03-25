package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jimmerzeel/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type responseBody struct {
		Chirp
	}

	decoder := json.NewDecoder(r.Body)
	input := requestBody{}
	err := decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request")
		return
	}
	dbParams := database.CreateChirpParams(input)
	dbChirp, err := cfg.db.CreateChirp(r.Context(), dbParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot create chirp")
		return
	}
	outputMsg, err := validateChirp(w, dbChirp.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp := responseBody{
		Chirp: Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      outputMsg,
			UserID:    dbChirp.UserID,
		},
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}

func validateChirp(w http.ResponseWriter, body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", fmt.Errorf("Chirp is too long")
	}
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	inputMsg := strings.Split(body, " ")
	for i, word := range inputMsg {
		if slices.Contains(profanity, strings.ToLower(word)) {
			inputMsg[i] = "****"
		}
	}
	return strings.Join(inputMsg, " "), nil
}

// func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
// 	type requestBody struct {
// 		Body string `json:"body"`
// 	}
// 	type responseBody struct {
// 		CleanedBody string `json:"cleaned_body"`
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	chirp := requestBody{}
// 	err := decoder.Decode(&chirp)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Error decoding the incoming Chirp")
// 		return
// 	}
// 	const maxChirpLength = 140
// 	if len(chirp.Body) > maxChirpLength {
// 		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
// 		return
// 	}
// 	profanity := []string{"kerfuffle", "sharbert", "fornax"}
// 	inputMsg := strings.Split(chirp.Body, " ")
// 	for i, word := range inputMsg {
// 		if slices.Contains(profanity, strings.ToLower(word)) {
// 			inputMsg[i] = "****"
// 		}
// 	}
// 	outputMsg := strings.Join(inputMsg, " ")
// 	respondWithJSON(w, http.StatusOK, responseBody{CleanedBody: outputMsg})
// }
