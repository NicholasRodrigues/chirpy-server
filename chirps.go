package main

import (
	database2 "github.com/NicholasRodrigues/chirpy-server/internal/database"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) insertChirpHandler(w http.ResponseWriter, r *http.Request) {
	var params database2.Chirp

	if err := decodeRequestBody(r, &params); err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid request body")
		return
	}

	db, err := database2.NewDB(cfg.dbPath)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error creating database")
		return
	}

	chirp, err := db.CreateChirp(params.Message)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error creating chirp")
		return
	}

	sendJSONResponse(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database2.NewDB(cfg.dbPath)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error creating database")
		return
	}

	chirps, err := db.GetChirps()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error getting chirps")
		return
	}

	sendJSONResponse(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database2.NewDB(cfg.dbPath)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error creating database")
		return
	}

	id := r.PathValue("id")

	intId, err := strconv.Atoi(id)

	chirp, err := db.GetChirpById(intId)
	if err != nil {
		if err.Error() == "chirp not found" {
			handleError(w, err, http.StatusNotFound, "Chirp not found")
		} else {
			handleError(w, err, http.StatusInternalServerError, "Error getting chirp")
		}
		return
	}

	sendJSONResponse(w, http.StatusOK, chirp)
}
