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
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}

	if err := decodeRequestBody(r, &params); err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid request body")
		return
	}

	userResponse, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error creating user")
		return
	}

	sendJSONResponse(w, http.StatusCreated, userResponse)
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

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int64 `json:"expires_in_seconds"`
	}
	params := parameters{}

	if err := decodeRequestBody(r, &params); err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid request body")
		return
	}

	userResponse, err := cfg.DB.LoginUser(params.Email, params.Password)
	if err != nil {
		handleError(w, err, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	var expires int64
	if params.ExpiresInSeconds != nil {
		expires = *params.ExpiresInSeconds
	} else {
		expires = 0
	}

	jwtToken, err := cfg.createJWT(userResponse.ID, expires)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error creating JWT")
		return
	}

	userResponse.Token = jwtToken

	sendJSONResponse(w, http.StatusOK, userResponse)
}

// Update user handler must infer the user ID from the JWT token
func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}

	if err := decodeRequestBody(r, &params); err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get the user ID from the JWT token
	ctxUserID := r.Context().Value("userID")
	userID, err := strconv.Atoi(ctxUserID.(string))
	if err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userResponse, err := cfg.DB.UpdateUser(userID, params.Email, params.Password)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error updating user")
		return
	}

	sendJSONResponse(w, http.StatusOK, userResponse)
}
