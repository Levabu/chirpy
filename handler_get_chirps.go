package main

import (
	"net/http"
	"slices"

	"github.com/Levabu/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	authorID, _ := uuid.Parse(s)

	order := r.URL.Query().Get("sort")
	if !slices.Contains([]string{"asc", "desc"}, order) {
		order = "asc"
	}

	var chirps []database.Chirp
	var err error
	if authorID != uuid.Nil {
		chirps, err = cfg.DB.GetChirpsByAuthorID(r.Context(), authorID)
	} else {
		chirps, err = cfg.DB.GetChirps(r.Context())
	}

	if order == "desc" {
		slices.Reverse(chirps)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps", err)
	}

	respondWithJSON(w, http.StatusOK, chirps)
}