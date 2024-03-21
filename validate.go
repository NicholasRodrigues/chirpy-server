package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type validateInputJson struct {
	Text string `json:"body"`
}

type validateOutputJson struct {
	IsValid bool `json:"valid"`
}

type invalidOutputJson struct {
	Error string `json:"error"`
}

// Helper function to decode JSON request body
func decodeRequestBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

// Helper function to send JSON response
func sendJSONResponse(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

// Helper function to handle errors and log them
func handleError(w http.ResponseWriter, err error, statusCode int, errMsg string) {
	log.Printf("Error: %s", err)
	if errMsg != "" {
		sendJSONResponse(w, statusCode, invalidOutputJson{Error: errMsg})
	} else {
		w.WriteHeader(statusCode)
	}
}

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

	sendJSONResponse(w, http.StatusOK, validateOutputJson{IsValid: true})
}
