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

	cleaned, err := validateChirp(input.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: input.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, responseBody{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	})
}

func validateChirp(body string) (string, error) {
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
