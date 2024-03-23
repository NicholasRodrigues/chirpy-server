package main

import "github.com/NicholasRodrigues/chirpy-server/internal/database"

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}
