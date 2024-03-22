package main

import (
	"net/http"
	"strconv"
)

type User struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}

func (cfg *apiConfig) insertUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	params := parameters{}

	if err := decodeRequestBody(r, &params); err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error creating user")
		return
	}

	sendJSONResponse(w, http.StatusCreated, user)
}

func (cfg *apiConfig) getUserHandler(w http.ResponseWriter, r *http.Request) {
	stringID := r.PathValue("userID")
	userID, err := strconv.Atoi(stringID)
	if err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid user ID")
		return
	}

	dbUser, err := cfg.DB.GetUserById(userID)
	if err != nil {
		handleError(w, err, http.StatusNotFound, "User not found")
		return
	}

	user := User{
		ID:    dbUser.ID,
		Email: dbUser.Email,
	}

	sendJSONResponse(w, http.StatusOK, user)
}
