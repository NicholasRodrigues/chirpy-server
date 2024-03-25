package main

import (
	"encoding/json"
	"github.com/NicholasRodrigues/chirpy-server/internal/auth"
	"net/http"
	"strconv"
	"time"
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
		Email    string `json:"email"`
		Password string `json:"password"`
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

	jwtToken, err := auth.MakeJWT(
		userResponse.ID,
		cfg.jwtSecret,
		time.Hour,
		auth.TokenTypeAccess,
	)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error creating JWT")
		return
	}

	refreshToken, err := auth.MakeJWT(
		userResponse.ID,
		cfg.jwtSecret,
		time.Hour*24*30*6,
		auth.TokenTypeRefresh,
	)

	userResponse.AccessToken = jwtToken
	userResponse.RefreshToken = refreshToken

	sendJSONResponse(w, http.StatusOK, userResponse)
}

// Update user handler must infer the user ID from the JWT token
//func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
//	type parameters struct {
//		Email    string `json:"email"`
//		Password string `json:"password"`
//	}
//	params := parameters{}
//
//	if err := decodeRequestBody(r, &params); err != nil {
//		handleError(w, err, http.StatusBadRequest, "Invalid request body")
//		return
//	}
//
//	// Get the user ID from the JWT token
//	ctxUserID := r.Context().Value("userID")
//	userID, err := strconv.Atoi(ctxUserID.(string))
//	if err != nil {
//		handleError(w, err, http.StatusBadRequest, "Invalid user ID")
//		return
//	}
//
//	userResponse, err := cfg.DB.UpdateUser(userID, params.Email, params.Password)
//	if err != nil {
//		handleError(w, err, http.StatusInternalServerError, "Error updating user")
//		return
//	}

//	sendJSONResponse(w, http.StatusOK, userResponse)
//}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		handleError(w, err, http.StatusUnauthorized, "Invalid token")
		return
	}
	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		handleError(w, err, http.StatusUnauthorized, "Invalid token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid request body")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error hashing password")
		return
	}

	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := cfg.DB.UpdateUser(userIDInt, params.Email, hashedPassword)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error updating user")
		return
	}

	sendJSONResponse(w, http.StatusOK, response{User: User{Email: user.Email, ID: user.ID}})

}
