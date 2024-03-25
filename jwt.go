package main

import (
	"context"
	"fmt"
	"github.com/NicholasRodrigues/chirpy-server/internal/auth"
	"net/http"
)

// This function is a middleware that checks for a JWT token in the Authorization header of the request.
// If the token is valid, the request is passed to the next handler. If the token is invalid, the middleware
// returns a 401 Unauthorized response.
func (cfg *apiConfig) middlewareJWT(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.GetBearerToken(r.Header)

		if tokenString == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		// Parse the token
		userId, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)

		if err != nil {
			fmt.Println("Error parsing token: ", err)
			http.Error(w, "Unauthorized 2", http.StatusUnauthorized)
			return
		}

		// The token is valid and is an access token, continue processing the request
		ctx := context.WithValue(r.Context(), "userID", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid token")
		return
	}

	isRevoked, err := cfg.DB.IsTokenRevoked(refreshToken)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error checking token")
		return
	}
	if isRevoked {
		handleError(w, err, http.StatusUnauthorized, "Token revoked")
		return
	}

	accessToken, err := auth.RefreshToken(refreshToken, cfg.jwtSecret)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error refreshing token")
		return
	}

	sendJSONResponse(w, http.StatusOK, response{Token: accessToken})
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		handleError(w, err, http.StatusBadRequest, "Invalid token")
		return
	}

	err = cfg.DB.RevokeToken(refreshToken)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError, "Error revoking token")
		return
	}

	sendJSONResponse(w, http.StatusOK, nil)
}
