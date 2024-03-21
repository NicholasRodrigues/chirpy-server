package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Define the struct for input JSON
type validateInputJson struct {
	Text string `json:"body"`
}

// Updated struct for output JSON
type validateOutputJson struct {
	CleanedBody string `json:"cleaned_body"`
}

// Function to decode JSON request body
func decodeRequestBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

// Function to send JSON response
func sendJSONResponse(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}

// Function to handle errors and log them
func handleError(w http.ResponseWriter, err error, statusCode int, errMsg string) {
	log.Printf("Error: %s", err)
	if errMsg != "" {
		sendJSONResponse(w, statusCode, map[string]string{"error": errMsg})
	} else {
		w.WriteHeader(statusCode)
	}
}

// Profanity filter
func sanitizeText(text string) string {
	profaneWords := map[string]string{
		"kerfuffle": "****",
		"sharbert":  "****",
		"fornax":    "****",
	}

	words := strings.Fields(text)
	for i, word := range words {
		lowerWord := strings.ToLower(word) // Convert word to lowercase for comparison
		if replacement, exists := profaneWords[lowerWord]; exists {
			words[i] = replacement
		}
	}

	return strings.Join(words, " ")
}

// Updated handler with profanity filtering
func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	var params validateInputJson
	if err := decodeRequestBody(r, &params); err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error decoding parameters")
		return
	}

	if len(params.Text) > 140 {
		handleError(w, nil, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedText := sanitizeText(params.Text)
	sendJSONResponse(w, http.StatusOK, validateOutputJson{CleanedBody: cleanedText})
}
