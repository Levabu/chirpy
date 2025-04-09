package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Levabu/chirpy/internal/auth"
	"github.com/Levabu/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type returnParameters struct {
		ID             uuid.UUID `json:"id"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		Email          string    `json:"email"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding json", err)
		return
	}

	// check auth
	access_token, _ := auth.GetBearerToken(r.Header)
	userID, err := auth.ValidateJWT(access_token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no valid authorization was provided", nil)
		return
	}
	user, err := cfg.DB.GetUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no valid authorization was provided", nil)
		return
	}

	if len(params.Password) < 1 {
		respondWithError(w, http.StatusBadRequest, "password must be at least 6 characters long", nil)
		return
	}
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}

	updatedUser, err := cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: user.ID,
		Email: params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnParameters{
		ID: updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email: updatedUser.Email,
	})
}
