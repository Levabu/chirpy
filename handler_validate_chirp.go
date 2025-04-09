package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding json", err)
		return
	}

	if isValid, err := isValidChirp(params.Body); !isValid {
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: replaceBadWords(params.Body),
	})
}

func isValidChirp(s string) (bool, error) {
	const maxChirpLength = 140
	if len(s) > maxChirpLength {
		return false, fmt.Errorf("chirp is too long")
	}
	return true, nil
}

func replaceBadWords(s string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(s, " ")
	for i, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
