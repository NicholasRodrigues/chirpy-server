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

func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := validateInputJson{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Text) > 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(invalidOutputJson{Error: "Chirp is too long"})
		return
	}

	respBody := validateOutputJson{
		IsValid: true,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
