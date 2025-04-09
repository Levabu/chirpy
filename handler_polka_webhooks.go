package main

import (
	"encoding/json"
	"net/http"

	"github.com/Levabu/chirpy/internal/auth"
	"github.com/google/uuid"
)

const (
	EVENT_USER_UPGRADED = "user.upgraded"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, _ := auth.GetApiKey(r.Header)
	if apiKey != cfg.PolkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding json", err)
		return
	}

	if params.Event != EVENT_USER_UPGRADED {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.DB.UpgradeUser(r.Context(), params.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}