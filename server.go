package main

import (
	"fmt"
	"github.com/NicholasRodrigues/chirpy-server/internal/database"
	"log"
	"net/http"
	"os"
)

func createServer() *http.Server {
	const port = "8080"
	const filepathRoot = "."
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	apiConfig := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	fsHandler := apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /api/reset", apiConfig.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", apiConfig.validateChirpHandler)
	mux.HandleFunc("POST /api/chirps", apiConfig.insertChirpHandler)
	mux.HandleFunc("GET /api/chirps", apiConfig.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.getChirpHandler)
	mux.HandleFunc("GET /api/users/{userID}", apiConfig.getUserHandler)
	mux.HandleFunc("POST /api/users", apiConfig.insertUserHandler)
	mux.HandleFunc("POST /api/login", apiConfig.loginUserHandler)

	//mux.Handle("PUT /api/users", apiConfig.middlewareJWT(http.HandlerFunc(apiConfig.updateUserHandler)))
	mux.HandleFunc("PUT /api/users", apiConfig.updateUserHandler)

	mux.HandleFunc("POST /api/refresh", apiConfig.refreshHandler)
	mux.HandleFunc("POST /api/revoke", apiConfig.revokeHandler)

	mux.HandleFunc("GET /admin/metrics", apiConfig.metricsHandler)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fmt.Println("Server created at: " + "http://localhost:" + port + "/app/")
	return server
}
