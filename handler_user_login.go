package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Levabu/chirpy/internal/auth"
	"github.com/Levabu/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type returnParameters struct {
		ID             uuid.UUID `json:"id"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		Email          string    `json:"email"`
		Token 				 string    `json:"token"`
		RefreshToken 	 string    `json:"refresh_token"`
		IsChirpyRed    bool      `json:"is_chirpy_red"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding json", err)
		return
	}


	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating refresh token", err)
		return
	}
	err = cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.JWTSecret, auth.TOKEN_EXP_TIME)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error signing jwt token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnParameters{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: accessToken,
		RefreshToken: refreshToken,
		IsChirpyRed: user.IsChirpyRed,
	})
}