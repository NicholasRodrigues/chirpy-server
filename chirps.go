package main

import (
	"errors"
	"github.com/NicholasRodrigues/chirpy-server/internal/auth"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Chirp struct {
	AuthorId int    `json:"author_id"`
	Message  string `json:"body"`
	ID       int    `json:"id"`
}

func (cfg *apiConfig) insertChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Message string `json:"body"`
	}
	chirp := parameters{}
	if err := decodeRequestBody(r, &chirp); err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid token")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		handleError(w, err, http.StatusUnauthorized, "Invalid token")
		return
	}

	cleanedBody, err := validateChirp(chirp.Message)
	if err != nil {
		handleError(w, err, http.StatusBadRequest, "Chirp is too long")
		return
	}

	dbChirp, err := cfg.DB.CreateChirp(cleanedBody, subject)
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
		AuthorId: dbChirp.AuthorId,
		ID:       dbChirp.ID,
		Message:  dbChirp.Message,
	}

	sendJSONResponse(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	authorIDParam := r.URL.Query().Get("author_id")

	if authorIDParam != "" {
		authorID, err := strconv.Atoi(authorIDParam)
		if err != nil {
			handleError(w, err, http.StatusBadRequest, "Invalid author ID")
			return
		}

		dbChirps, err := cfg.DB.GetChirpsByAuthorId(authorID)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, "Error retrieving chirps")
			return
		}

		chirps := []Chirp{}
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, Chirp{
				AuthorId: dbChirp.AuthorId,
				ID:       dbChirp.ID,
				Message:  dbChirp.Message,
			})
		}

		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})

		sendJSONResponse(w, http.StatusOK, chirps)
		return
	}

	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error retrieving chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			AuthorId: dbChirp.AuthorId,
			ID:       dbChirp.ID,
			Message:  dbChirp.Message,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	sendJSONResponse(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	chirp, err := cfg.DB.GetChirpById(chirpID)
	if err != nil {
		handleError(w, err, http.StatusNotFound, "Chirp not found")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid token")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		handleError(w, err, http.StatusUnauthorized, "Invalid token")
		return
	}

	authorId := strconv.Itoa(chirp.AuthorId)

	if authorId != subject {
		handleError(w, err, http.StatusForbidden, "Unauthorized")
		return
	}

	err = cfg.DB.DeleteChirp(chirpID)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error deleting chirp")
		return
	}

	sendJSONResponse(w, http.StatusOK, nil)
}
