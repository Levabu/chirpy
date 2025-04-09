package main

import (
	"net/http"
	"time"

	"github.com/Levabu/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type returnParameters struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No authorization was provided", nil)
		return
	}

	rt, err := cfg.DB.GetRefreshToken(r.Context(), refreshToken)
	if err != nil ||
	  rt.ExpiresAt.Before(time.Now().UTC()) || 
		rt.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "No authorization was provided", err)
		return
	}

	accessToken, err := auth.MakeJWT(rt.UserID, cfg.JWTSecret, auth.TOKEN_EXP_TIME)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error signing access token", err)
		return
	}


	respondWithJSON(w, http.StatusOK, returnParameters{
		Token: accessToken,
	})
}