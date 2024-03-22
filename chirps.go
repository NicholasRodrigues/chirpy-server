package main

import (
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Chirp struct {
	Message string `json:"body"`
	ID      int    `json:"id"`
}

func (cfg *apiConfig) insertChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirp Chirp
	if err := decodeRequestBody(r, &chirp); err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid request body")
		return
	}

	cleanedBody, err := validateChirp(chirp.Message)
	if err != nil {
		handleError(w, err, http.StatusBadRequest, "Chirp is too long")
		return
	}

	dbChirp, err := cfg.DB.CreateChirp(cleanedBody)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error creating chirp")
		return
	}

	sendJSONResponse(w, http.StatusCreated, dbChirp)
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirpById(chirpID)
	if err != nil {
		handleError(w, err, http.StatusNotFound, "Chirp not found")
		return
	}

	chirp := Chirp{
		ID:      dbChirp.ID,
		Message: dbChirp.Message,
	}

	sendJSONResponse(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error retrieving chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:      dbChirp.ID,
			Message: dbChirp.Message,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	sendJSONResponse(w, http.StatusOK, chirps)
}
