package main

import (
	database2 "github.com/NicholasRodrigues/chirpy-server/internal/database"
	"net/http"
)

func (cfg *apiConfig) insertUserHandler(w http.ResponseWriter, r *http.Request) {
	var params database2.User

	if err := decodeRequestBody(r, &params); err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid request body")
		return
	}

	db, err := database2.NewDB(cfg.dbPath)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error creating database")
		return
	}

	user, err := db.CreateUser(params.Email)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error creating user")
		return
	}

	sendJSONResponse(w, http.StatusCreated, user)
}

func (cfg *apiConfig) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database2.NewDB(cfg.dbPath)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error creating database")
		return
	}

	users, err := db.GetUsers()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error getting users")
		return
	}

	sendJSONResponse(w, http.StatusOK, users)
}
