package main

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"time"
)

func (cfg *apiConfig) createJWT(userID int, expiresInSeconds int64) (string, error) {
	if expiresInSeconds == 0 || expiresInSeconds > 86400 {
		expiresInSeconds = 86400
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresInSeconds) * time.Second)),
		Subject:   strconv.Itoa(userID),
	})

	token, err := claims.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

// This function is a middleware that checks for a JWT token in the Authorization header of the request.
// If the token is valid, the request is passed to the next handler. If the token is invalid, the middleware
// returns a 401 Unauthorized response.
func (cfg *apiConfig) middlewareJWT(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := authorizationHeader[7:]
		claims := jwt.RegisteredClaims{}
		_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.jwtSecret), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
