package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
