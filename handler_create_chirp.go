package main

import (
	"encoding/json"
	"net/http"

	"github.com/Levabu/chirpy/internal/auth"
	"github.com/Levabu/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding json", err)
		return
	}

	token, _ := auth.GetBearerToken(r.Header)
	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no valid authorization was provided", nil)
		return
	}

	if isValid, err := isValidChirp(params.Body); !isValid {
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	params.Body = replaceBadWords(params.Body)

	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}
