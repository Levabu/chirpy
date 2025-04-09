package main

import (
	"net/http"

	"github.com/Levabu/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No valid chirp id was provided", nil)
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

	chirp, err := cfg.DB.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No chirp found", err)
		return
	}

	if user.ID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Not authorized", err)
		return
	}

	err = cfg.DB.DeleteChirp(r.Context(), chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No chirp found", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}