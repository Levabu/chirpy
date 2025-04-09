package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Levabu/chirpy/internal/auth"
	"github.com/Levabu/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type returnParameters struct {
		ID             uuid.UUID `json:"id"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		Email          string    `json:"email"`
		IsChirpyRed    bool      `json:"is_chirpy_red"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding json", err)
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

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			respondWithError(w, http.StatusConflict, "user with that email already exists", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnParameters{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
