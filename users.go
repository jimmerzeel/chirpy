package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jimmerzeel/chirpy/internal/auth"
	"github.com/jimmerzeel/chirpy/internal/database"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type responseBody struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	input := requestBody{}
	err := decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request")
		return
	}
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error handling user input")
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          input.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot create user")
		return
	}
	user := responseBody{
		User: User{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt,
			UpdatedAt:   dbUser.UpdatedAt,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type responseBody struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	input := requestBody{}
	err := decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request")
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(dbUser.HashedPassword, input.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token")
		return
	}
	refreshExpiration := time.Now().AddDate(0, 0, 60)

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: refreshExpiration,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error storing refresh token in the database")
		return
	}

	user := responseBody{
		User: User{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt,
			UpdatedAt:   dbUser.UpdatedAt,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type responseBody struct {
		User
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Cannot find access token")
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	input := requestBody{}
	err = decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request")
		return
	}

	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error handling user input")
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          input.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot update user")
		return
	}
	respondWithJSON(w, http.StatusOK, responseBody{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		}})
}
